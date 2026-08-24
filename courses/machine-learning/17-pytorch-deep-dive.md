# Chapter 17: PyTorch — The Deep Learning Framework

> **"PyTorch treats deep learning like regular Python programming. The computation graph is built as you run code — not before. This 'define by run' approach changed how researchers think about and debug neural networks."**

---

## Table of Contents
1. [Why PyTorch](#1-why-pytorch)
2. [Tensors — The Fundamental Data Structure](#2-tensors)
3. [Autograd — Automatic Differentiation](#3-autograd--automatic-differentiation)
4. [nn.Module — Building Networks](#4-nnmodule--building-networks)
5. [Built-In Layers and Functions](#5-built-in-layers-and-functions)
6. [Dataset and DataLoader](#6-dataset-and-dataloader)
7. [The Complete Training Loop](#7-the-complete-training-loop)
8. [GPU Utilization](#8-gpu-utilization)
9. [Model Saving and Loading](#9-model-saving-and-loading)
10. [Mixed Precision Training](#10-mixed-precision-training)
11. [Debugging PyTorch](#11-debugging-pytorch)
12. [Summary and What's Next](#12-summary)

---

## 1. Why PyTorch

PyTorch was released by Facebook AI Research (FAIR) in 2016. By 2020, it had become the dominant framework in academic research, and by 2023, it was widespread in production systems too.

### Dynamic vs Static Computation Graphs

**TensorFlow 1.x (static graph — "define then run"):**
```python
# TF1 style — build graph first, then run
graph = tf.Graph()
with graph.as_default():
    x = tf.placeholder(tf.float32, shape=[None, 784])
    W = tf.Variable(tf.zeros([784, 10]))
    y = tf.matmul(x, W)
    # ... many more lines of graph construction

with tf.Session(graph=graph) as sess:
    result = sess.run(y, feed_dict={x: data})
```

**PyTorch (dynamic graph — "define by run"):**
```python
# PyTorch style — just write regular Python
x = torch.tensor(data)
W = torch.zeros(784, 10, requires_grad=True)
y = x @ W   # graph is built HERE, as this line executes
```

**Why dynamic graphs matter:**
- Use regular Python control flow (if/else, for loops, recursion)
- Inspect tensors at any point with print()
- Change network behavior based on data (variable-length sequences, tree structures)
- Debugging is standard Python debugging

### PyTorch vs TensorFlow: Honest Comparison

| | PyTorch | TensorFlow 2.x |
|-|---------|---------------|
| Ease of use | Excellent — Pythonic | Good — mostly Pythonic now |
| Research dominance | ~75% of papers (2024) | ~25% |
| Production ecosystem | TorchScript, ONNX, TorchServe | TF Serving, TF Lite, TF.js |
| Mobile/Edge | TorchScript, ExecuTorch | TF Lite (more mature) |
| TPU support | Partial (via XLA) | Excellent (native) |
| Ecosystem | PyPI, HuggingFace, Lightning | Keras, TF Hub |
| Default in academia | Yes | Was, less so now |

**Verdict**: For learning, research, and most production use cases — use PyTorch. If you're deploying to Android/iOS at scale, TF Lite is worth considering. If you're using Google TPUs, TensorFlow is better supported.

### Installation

```bash
# CPU only (for learning, no GPU required):
pip install torch torchvision torchaudio

# CUDA 12.1 (NVIDIA GPU):
pip install torch torchvision torchaudio --index-url https://download.pytorch.org/whl/cu121

# Mac M1/M2 (Metal Performance Shaders):
pip install torch torchvision torchaudio

# Verify:
python -c "import torch; print(torch.__version__); print(torch.cuda.is_available())"
```

---

## 2. Tensors

A tensor is an n-dimensional array. It is to PyTorch what ndarray is to NumPy — the fundamental data structure. The key difference: tensors can live on GPUs and can automatically track gradients.

### Creating Tensors

```python
import torch
import numpy as np

# ─── From Python data ────────────────────────────────────────────
t1 = torch.tensor([1.0, 2.0, 3.0])           # 1D tensor from list
t2 = torch.tensor([[1, 2], [3, 4]])           # 2D tensor (matrix)
t3 = torch.tensor(3.14)                       # 0D tensor (scalar)

# ─── Factory functions ───────────────────────────────────────────
zeros = torch.zeros(3, 4)                     # 3×4 tensor of zeros
ones  = torch.ones(2, 3, 4)                   # 3D tensor of ones
empty = torch.empty(5, 5)                     # uninitialized (garbage values)
rand  = torch.rand(3, 3)                      # uniform [0,1)
randn = torch.randn(3, 3)                     # standard normal N(0,1)
arange = torch.arange(0, 10, 2)              # [0, 2, 4, 6, 8]
linspace = torch.linspace(0, 1, 5)           # [0.00, 0.25, 0.50, 0.75, 1.00]

# Like-sized (copy shape and dtype, not values):
zeros_like = torch.zeros_like(t2)            # same shape as t2
ones_like  = torch.ones_like(t2)

# ─── From NumPy (SHARED MEMORY — changes to one affect the other!) ──
np_array = np.array([1.0, 2.0, 3.0])
torch_from_np = torch.from_numpy(np_array)   # shares memory with np_array!
torch_copy = torch.tensor(np_array)          # makes a COPY (safe)

# NumPy from tensor:
np_from_torch = torch_from_np.numpy()        # also shares memory!
safe_np = torch_from_np.detach().numpy()     # safe copy (also detaches grad)

print(f"Shared memory demo:")
np_array[0] = 99
print(f"np_array[0] = {np_array[0]}")
print(f"torch_from_np[0] = {torch_from_np[0]}")  # also 99!
```

### Data Types

```python
# Float types:
f32 = torch.zeros(3, dtype=torch.float32)   # default, most common
f64 = torch.zeros(3, dtype=torch.float64)   # double precision
f16 = torch.zeros(3, dtype=torch.float16)   # half precision (GPU training)
bf16 = torch.zeros(3, dtype=torch.bfloat16) # brain float (modern GPUs)

# Integer types:
i64 = torch.zeros(3, dtype=torch.int64)     # long (class indices in CrossEntropyLoss)
i32 = torch.zeros(3, dtype=torch.int32)
i8  = torch.zeros(3, dtype=torch.int8)

# Boolean:
bools = torch.tensor([True, False, True], dtype=torch.bool)

# Convert dtype:
x = torch.randn(3)
x_double = x.double()    # to float64
x_int = x.int()          # to int32
x_long = x.long()        # to int64
x_half = x.half()        # to float16

# Check type:
print(x.dtype)            # torch.float32
print(x.shape)            # torch.Size([3])
print(x.ndim)             # 1
print(x.numel())          # 3  (total elements)
```

### Tensor Operations

```python
import torch

a = torch.tensor([[1.0, 2.0], [3.0, 4.0]])
b = torch.tensor([[5.0, 6.0], [7.0, 8.0]])

# Arithmetic (element-wise):
c = a + b          # or torch.add(a, b)
c = a - b          # subtraction
c = a * b          # element-wise multiply
c = a / b          # element-wise divide
c = a ** 2         # element-wise power

# Matrix multiplication:
c = a @ b          # matmul (preferred)
c = torch.mm(a, b) # 2D only
c = torch.matmul(a, b)  # general, works for batched too

# Comparison:
mask = a > 2       # tensor([[False, False], [ True,  True]])

# Reductions:
a.sum()            # sum of all elements
a.sum(dim=0)       # sum along rows (result: 1D, size 2)
a.sum(dim=1)       # sum along columns
a.mean()
a.max()
a.max(dim=1)       # returns (values, indices) tuple
a.argmax()         # index of maximum value
a.min()

# Shape operations:
x = torch.randn(2, 3, 4)
x.shape             # torch.Size([2, 3, 4])
x.reshape(2, 12)    # reshape (copies if needed)
x.view(2, 12)       # reshape (never copies — must be contiguous)
x.contiguous()      # make contiguous in memory
x.permute(0, 2, 1)  # reorder dimensions: (2,3,4) → (2,4,3)
x.transpose(1, 2)   # swap two dimensions: equivalent to permute(0,2,1)

# Adding/removing dimensions:
y = torch.randn(3)
y.unsqueeze(0)      # (3,) → (1, 3)
y.unsqueeze(1)      # (3,) → (3, 1)
z = torch.randn(1, 3)
z.squeeze(0)        # (1, 3) → (3,)
z.squeeze()         # remove ALL size-1 dimensions

# Concatenation and stacking:
p = torch.randn(2, 3)
q = torch.randn(2, 3)
torch.cat([p, q], dim=0)    # (4, 3) — concatenate along existing dim
torch.cat([p, q], dim=1)    # (2, 6) — concatenate along dim 1
torch.stack([p, q], dim=0)  # (2, 2, 3) — NEW dimension created

# Indexing (just like NumPy):
x = torch.randn(5, 4)
x[0]               # first row
x[:, 0]            # first column
x[1:3, 2:]         # rows 1-2, columns 2-3
x[x > 0]           # boolean indexing — returns 1D tensor of all positive values

# .item() — extract scalar Python value:
loss = torch.tensor(0.4532)
python_float = loss.item()  # 0.4532 (Python float, not tensor)
# Use .item() when you want to print/log a loss value
```

### In-Place Operations

```python
x = torch.randn(3, 3)

# In-place operations have underscore suffix:
x.add_(1.0)        # x = x + 1.0  (in-place)
x.mul_(2.0)        # x = x * 2.0  (in-place)
x.zero_()          # x = 0         (in-place)
x.fill_(3.14)      # fill with value (in-place)

# WARNING: in-place operations on tensors that require grad will cause errors
# Avoid in-place ops on tensors in the computation graph
```

---

## 3. Autograd — Automatic Differentiation

Autograd is what makes PyTorch a deep learning framework rather than just a tensor library. It automatically computes gradients of any computation you perform on tensors.

### requires_grad and backward()

```python
import torch

# ─── Simple scalar example ────────────────────────────────────────
x = torch.tensor(3.0, requires_grad=True)  # tell PyTorch to track this

# Forward pass — builds computation graph
y = x ** 2 + 2 * x + 1   # y = (x+1)²

print(f"y = {y}")          # tensor(16., grad_fn=<AddBackward0>)
# Note: y has grad_fn, which is how PyTorch knows how to differentiate it

# Backward pass — compute gradients
y.backward()

# Gradient: dy/dx = 2x + 2 = 2(3) + 2 = 8
print(f"dy/dx at x=3: {x.grad}")   # tensor(8.)
```

### Computational Graph

```python
import torch

# Build a simple computation graph:
#   z = (x + y) * (x - y)

x = torch.tensor(2.0, requires_grad=True)
y = torch.tensor(3.0, requires_grad=True)

a = x + y          # a = 5
b = x - y          # b = -1
z = a * b          # z = -5

# Compute gradients
z.backward()

print(f"x.grad = {x.grad}")   # dz/dx = (x-y) + (x+y) = 2x = 4.0
print(f"y.grad = {y.grad}")   # dz/dy = (x+y)*(-1) + (x-y)*(1) = -2y = -6.0
```

### torch.no_grad() — Disabling Gradient Tracking

```python
import torch

model = SomeModel()   # assume defined elsewhere

x = torch.randn(10, 784)

# During inference: no need to track gradients → saves memory and computation
with torch.no_grad():
    output = model(x)
    # No computation graph is built here

# Equivalent decorator:
@torch.no_grad()
def predict(model, x):
    return model(x)

# Check if grad is tracked:
print(output.requires_grad)   # False (no_grad context disabled it)
```

### Gradient Accumulation Warning

```python
# COMMON MISTAKE: gradients accumulate by default!
# PyTorch ADDS to .grad rather than replacing it.

x = torch.tensor(3.0, requires_grad=True)

y = x ** 2
y.backward()
print(f"First backward: x.grad = {x.grad}")   # 6.0

y = x ** 2
y.backward()
print(f"Second backward: x.grad = {x.grad}")  # 12.0 !!!! (accumulated, not 6.0)

# Solution: call x.grad.zero_() or optimizer.zero_grad() before each backward
x.grad.zero_()
y = x ** 2
y.backward()
print(f"After zero: x.grad = {x.grad}")        # 6.0 (correct)
```

### Detach — Stop Gradient Flow

```python
# detach() creates a new tensor that shares data but has no grad_fn
# Useful when you want to use a value without including it in the computation graph

x = torch.randn(5, requires_grad=True)
y = x * 2

# y is part of the graph; y_detached is NOT
y_detached = y.detach()

# Use case: GAN training — detach fake image before passing to discriminator
# Use case: Computing metrics without computing gradients

# Also useful for plotting — matplotlib cannot handle grad tensors
loss_value = loss.detach().cpu().numpy()   # detach, move to CPU, convert to numpy
```

---

## 4. nn.Module — Building Networks

`nn.Module` is the base class for all neural network models and layers in PyTorch.

### Basic Structure

```python
import torch
import torch.nn as nn

class MyLinear(nn.Module):
    """A custom fully-connected layer implemented from scratch."""
    
    def __init__(self, in_features, out_features):
        super().__init__()   # ALWAYS call super().__init__() first
        
        # nn.Parameter: a tensor that is automatically included in
        # model.parameters() and thus updated by the optimizer
        self.weight = nn.Parameter(
            torch.randn(out_features, in_features) * 0.01
        )
        self.bias = nn.Parameter(
            torch.zeros(out_features)
        )
    
    def forward(self, x):
        # x shape: (batch_size, in_features)
        return x @ self.weight.T + self.bias


class TwoLayerMLP(nn.Module):
    """A two-hidden-layer MLP."""
    
    def __init__(self, in_dim, hidden_dim, out_dim):
        super().__init__()
        
        # Register submodules — PyTorch tracks their parameters automatically
        self.layer1 = nn.Linear(in_dim, hidden_dim)
        self.relu1 = nn.ReLU()
        self.layer2 = nn.Linear(hidden_dim, hidden_dim)
        self.relu2 = nn.ReLU()
        self.output = nn.Linear(hidden_dim, out_dim)
        
        # Custom initialization
        self._init_weights()
    
    def _init_weights(self):
        for module in self.modules():
            if isinstance(module, nn.Linear):
                nn.init.kaiming_normal_(module.weight, nonlinearity='relu')
                nn.init.zeros_(module.bias)
    
    def forward(self, x):
        # x shape: (batch_size, in_dim)
        x = self.relu1(self.layer1(x))
        x = self.relu2(self.layer2(x))
        x = self.output(x)
        return x    # raw logits, no activation here


# Instantiate and inspect
model = TwoLayerMLP(784, 512, 10)

# Count parameters
total_params = sum(p.numel() for p in model.parameters())
trainable_params = sum(p.numel() for p in model.parameters() if p.requires_grad)
print(f"Total params: {total_params:,}")
print(f"Trainable params: {trainable_params:,}")

# Iterate over named parameters
for name, param in model.named_parameters():
    print(f"  {name}: shape={param.shape}, requires_grad={param.requires_grad}")
```

### Training vs Evaluation Mode

```python
model = TwoLayerMLP(784, 512, 10)

# Training mode: Dropout uses random masks, BatchNorm uses batch statistics
model.train()

# Evaluation mode: Dropout disabled (all neurons active), BatchNorm uses running stats
model.eval()

# CRITICAL: Always call model.eval() during validation/inference!
# Forgetting this is one of the most common bugs in PyTorch code.
# Dropout during evaluation = stochastic predictions = wrong metrics
```

### nn.Sequential — Quick Model Building

```python
# nn.Sequential chains layers in order
model = nn.Sequential(
    nn.Linear(784, 512),
    nn.ReLU(),
    nn.Dropout(0.2),
    nn.Linear(512, 256),
    nn.ReLU(),
    nn.Dropout(0.2),
    nn.Linear(256, 10)
)

# Named Sequential (more readable):
from collections import OrderedDict

model = nn.Sequential(OrderedDict([
    ('fc1', nn.Linear(784, 512)),
    ('relu1', nn.ReLU()),
    ('drop1', nn.Dropout(0.2)),
    ('fc2', nn.Linear(512, 256)),
    ('relu2', nn.ReLU()),
    ('drop2', nn.Dropout(0.2)),
    ('output', nn.Linear(256, 10))
]))

# Access a specific layer:
model.fc1          # first linear layer
model[0]           # same thing
```

---

## 5. Built-In Layers and Functions

### Linear Layers

```python
import torch.nn as nn

# Fully connected:
fc = nn.Linear(in_features=784, out_features=256, bias=True)
# Weight shape: (256, 784)
# Bias shape: (256,)
# output = input @ weight.T + bias

# Multiple inputs:
x = torch.randn(32, 784)   # batch of 32 samples
out = fc(x)                # shape: (32, 256)
```

### Activation Functions

```python
# Functional API (stateless — can use as one-liner):
import torch.nn.functional as F

x = torch.randn(10)
F.relu(x)
F.sigmoid(x)
F.tanh(x)
F.gelu(x)
F.silu(x)           # SiLU = Swish
F.softmax(x, dim=0) # must specify dim

# Module API (has state — can be part of nn.Sequential):
relu    = nn.ReLU()
sigmoid = nn.Sigmoid()
tanh    = nn.Tanh()
gelu    = nn.GELU()
softmax = nn.Softmax(dim=1)  # dim=1 for (batch, classes)

# Softmax note: NEVER apply softmax before CrossEntropyLoss in PyTorch
# CrossEntropyLoss already includes LogSoftmax internally
# Applying softmax before it causes numerical instability and wrong results
```

### Dropout

```python
# Standard dropout (for fully connected layers):
dropout = nn.Dropout(p=0.5)     # p = probability of ZEROING a neuron

# 2D dropout (for convolutional feature maps — zeros entire feature maps):
dropout2d = nn.Dropout2d(p=0.1)

# How it works:
# Training mode: randomly zero p fraction of inputs, scale remaining by 1/(1-p)
# Eval mode: pass-through (no zeroing), no scaling needed

x = torch.randn(4, 8)
layer = nn.Dropout(0.5)

layer.train()
print(layer(x))   # ~50% of values are 0

layer.eval()
print(layer(x))   # all values unchanged (full output)
```

### Normalization Layers

```python
# BatchNorm for 1D data (fully connected):
# Normalizes across the BATCH dimension
bn1d = nn.BatchNorm1d(num_features=256)   # 256 = number of features/neurons

# BatchNorm for 2D data (conv layers):
# Normalizes across batch AND spatial dimensions (H×W)
bn2d = nn.BatchNorm2d(num_features=64)    # 64 = number of channels

# LayerNorm — normalizes across FEATURES (not batch)
# Used in Transformers; works with any batch size including 1
# Does NOT use running stats — same behavior in train and eval
ln = nn.LayerNorm(normalized_shape=256)    # normalize over last 256 dims
ln_seq = nn.LayerNorm([512, 256])          # normalize over last 2 dims

# GroupNorm — between BatchNorm and LayerNorm
gn = nn.GroupNorm(num_groups=8, num_channels=64)

# InstanceNorm — normalizes each sample independently (style transfer)
inn = nn.InstanceNorm2d(num_features=64)
```

### Embedding Layer

```python
# For text/token processing:
# Maps integer token IDs to dense vectors

vocab_size = 10000   # size of vocabulary
embed_dim = 256      # embedding dimension

embedding = nn.Embedding(vocab_size, embed_dim)
# Weight matrix: (10000, 256) — a lookup table

# Input: token IDs as integers
token_ids = torch.tensor([5, 23, 101, 7])   # shape: (4,) — a sequence of 4 tokens
embeddings = embedding(token_ids)            # shape: (4, 256)

# Batch of sequences:
batch = torch.tensor([[5, 23, 101], [7, 8, 9]])   # (batch=2, seq_len=3)
out = embedding(batch)                             # (2, 3, 256)

# Padding token (should not contribute to loss):
embedding_with_pad = nn.Embedding(vocab_size, embed_dim, padding_idx=0)
# Index 0 always maps to all-zeros vector, gradient is also 0 for it
```

### Convolutional Layers

```python
# 2D convolution (for images):
conv = nn.Conv2d(
    in_channels=3,      # RGB input
    out_channels=64,    # number of filters (output channels)
    kernel_size=3,      # 3×3 filter
    stride=1,           # step size
    padding=1           # 'same' padding to maintain spatial size
)
# Input shape:  (batch, 3, H, W)
# Output shape: (batch, 64, H, W)  (with padding=1, stride=1)

# 1D convolution (for sequences):
conv1d = nn.Conv1d(in_channels=256, out_channels=512, kernel_size=3, padding=1)

# Transposed convolution (upsampling):
convT = nn.ConvTranspose2d(64, 32, kernel_size=2, stride=2)  # doubles spatial dims

# Depthwise separable convolution (efficient):
# Step 1: depthwise (group=in_channels)
dw = nn.Conv2d(in_channels=64, out_channels=64, kernel_size=3,
               padding=1, groups=64)  # each channel convolved independently
# Step 2: pointwise (1×1 conv)
pw = nn.Conv2d(in_channels=64, out_channels=128, kernel_size=1)
```

### Pooling

```python
# Max pooling:
maxpool = nn.MaxPool2d(kernel_size=2, stride=2)   # halves spatial dimensions

# Average pooling:
avgpool = nn.AvgPool2d(kernel_size=2, stride=2)

# Global Average Pooling (reduce H×W → 1×1 per channel):
gap = nn.AdaptiveAvgPool2d(output_size=(1, 1))
# Input: (batch, C, H, W) → Output: (batch, C, 1, 1)
# Then: x.flatten(1) → (batch, C)

# Global Max Pooling:
gmp = nn.AdaptiveMaxPool2d(output_size=(1, 1))
```

---

## 6. Dataset and DataLoader

### Custom Dataset

```python
import os
import torch
from torch.utils.data import Dataset, DataLoader
from PIL import Image
import torchvision.transforms as transforms

class ImageFolderDataset(Dataset):
    """
    Loads images from a directory structure:
      root/
        class_a/
          img1.jpg
          img2.jpg
        class_b/
          img3.jpg
          ...
    """
    
    def __init__(self, root_dir, transform=None):
        self.root_dir = root_dir
        self.transform = transform
        
        # Build list of (path, label) pairs
        self.samples = []
        self.classes = sorted(os.listdir(root_dir))
        self.class_to_idx = {cls: idx for idx, cls in enumerate(self.classes)}
        
        for class_name in self.classes:
            class_dir = os.path.join(root_dir, class_name)
            if not os.path.isdir(class_dir):
                continue
            for img_name in os.listdir(class_dir):
                if img_name.lower().endswith(('.jpg', '.jpeg', '.png')):
                    path = os.path.join(class_dir, img_name)
                    label = self.class_to_idx[class_name]
                    self.samples.append((path, label))
    
    def __len__(self):
        """Required: return total number of samples."""
        return len(self.samples)
    
    def __getitem__(self, idx):
        """
        Required: return (sample, label) for a given index.
        This is called by DataLoader for each sample in a batch.
        """
        img_path, label = self.samples[idx]
        
        # Load image (PIL Image)
        image = Image.open(img_path).convert('RGB')
        
        # Apply transforms if specified
        if self.transform is not None:
            image = self.transform(image)
        
        return image, label
    
    def get_class_counts(self):
        """Analyze class distribution."""
        from collections import Counter
        labels = [label for _, label in self.samples]
        counts = Counter(labels)
        return {self.classes[k]: v for k, v in sorted(counts.items())}
```

### Transforms

```python
import torchvision.transforms as transforms

# Training transforms (with augmentation):
train_transforms = transforms.Compose([
    transforms.Resize(256),                           # resize shorter side to 256
    transforms.RandomCrop(224),                       # random 224×224 crop
    transforms.RandomHorizontalFlip(p=0.5),          # flip with 50% probability
    transforms.ColorJitter(brightness=0.2,            # random brightness/contrast
                           contrast=0.2,
                           saturation=0.2,
                           hue=0.1),
    transforms.ToTensor(),                            # PIL Image → float tensor [0,1]
    transforms.Normalize(
        mean=[0.485, 0.456, 0.406],                  # ImageNet mean
        std=[0.229, 0.224, 0.225]                    # ImageNet std
    )
])

# Validation/Test transforms (NO augmentation — deterministic):
val_transforms = transforms.Compose([
    transforms.Resize(256),
    transforms.CenterCrop(224),                      # center crop (not random)
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406],
                         std=[0.229, 0.224, 0.225])
])

# NOTE: Always normalize with the same mean/std used during training.
# For transfer learning from ImageNet, use ImageNet mean/std.
# For training from scratch: compute mean/std from YOUR dataset.
```

### DataLoader

```python
from torch.utils.data import DataLoader, random_split

# Create full dataset
full_dataset = ImageFolderDataset('./data/flowers', transform=train_transforms)

# Split into train/val
train_size = int(0.8 * len(full_dataset))
val_size = len(full_dataset) - train_size
train_dataset, val_dataset = random_split(full_dataset, [train_size, val_size])

# Override val dataset's transforms (val should not have augmentation)
# Note: this is a common issue with random_split — the transform is on the original dataset
# Better approach: create two separate Dataset instances with different transforms

# DataLoaders:
train_loader = DataLoader(
    train_dataset,
    batch_size=32,
    shuffle=True,            # shuffle training data each epoch
    num_workers=4,           # parallel data loading (use 0 on Windows/debug)
    pin_memory=True,         # speeds up GPU transfer (if using CUDA)
    drop_last=True           # drop last incomplete batch (stabilizes BatchNorm)
)

val_loader = DataLoader(
    val_dataset,
    batch_size=64,           # can use larger batch for validation (no gradients)
    shuffle=False,           # never shuffle validation
    num_workers=4,
    pin_memory=True
)

# Inspect a batch:
images, labels = next(iter(train_loader))
print(f"Image batch shape: {images.shape}")    # (32, 3, 224, 224)
print(f"Labels batch shape: {labels.shape}")   # (32,)
print(f"Image dtype: {images.dtype}")          # torch.float32
print(f"Labels dtype: {labels.dtype}")         # torch.int64

# Torchvision built-in datasets:
import torchvision.datasets as datasets

mnist = datasets.MNIST('./data', train=True, download=True, transform=transforms.ToTensor())
cifar10 = datasets.CIFAR10('./data', train=True, download=True, transform=train_transforms)
imagenet = datasets.ImageFolder('./data/imagenet/train', transform=train_transforms)
```

---

## 7. The Complete Training Loop

A production-quality training loop includes: train/val phases, checkpointing, LR scheduling, gradient clipping, and logging.

```python
import torch
import torch.nn as nn
import torch.optim as optim
from torch.utils.tensorboard import SummaryWriter
import time
import os

def train_model(
    model,
    train_loader,
    val_loader,
    n_epochs=30,
    learning_rate=1e-3,
    weight_decay=1e-4,
    checkpoint_dir='./checkpoints',
    log_dir='./runs',
    device='cuda'
):
    """
    Production-quality training loop.
    
    Returns: best_val_accuracy, training history dict
    """
    # ─── Setup ────────────────────────────────────────────────────
    os.makedirs(checkpoint_dir, exist_ok=True)
    
    model = model.to(device)
    
    criterion = nn.CrossEntropyLoss()
    optimizer = optim.AdamW(
        model.parameters(),
        lr=learning_rate,
        weight_decay=weight_decay
    )
    
    # Cosine annealing with warm restarts
    scheduler = optim.lr_scheduler.CosineAnnealingLR(
        optimizer,
        T_max=n_epochs,
        eta_min=1e-6
    )
    
    # TensorBoard writer
    writer = SummaryWriter(log_dir)
    
    # Training history
    history = {
        'train_loss': [], 'val_loss': [],
        'train_acc': [], 'val_acc': [],
        'lr': []
    }
    
    best_val_acc = 0.0
    best_epoch = 0
    
    # ─── Training loop ────────────────────────────────────────────
    for epoch in range(n_epochs):
        epoch_start = time.time()
        
        # ─── TRAINING PHASE ─────────────────────────────────────
        model.train()   # <-- critical
        train_loss = 0.0
        train_correct = 0
        train_total = 0
        
        for batch_idx, (inputs, targets) in enumerate(train_loader):
            inputs = inputs.to(device, non_blocking=True)
            targets = targets.to(device, non_blocking=True)
            
            # Zero gradients
            optimizer.zero_grad()
            
            # Forward pass
            outputs = model(inputs)
            loss = criterion(outputs, targets)
            
            # Backward pass
            loss.backward()
            
            # Gradient clipping (prevents explosion, especially useful for RNNs)
            torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
            
            # Weight update
            optimizer.step()
            
            # Accumulate stats
            train_loss += loss.item() * inputs.size(0)
            _, predicted = outputs.max(1)
            train_total += targets.size(0)
            train_correct += predicted.eq(targets).sum().item()
        
        train_loss = train_loss / train_total
        train_acc = 100.0 * train_correct / train_total
        
        # ─── VALIDATION PHASE ───────────────────────────────────
        model.eval()   # <-- critical
        val_loss = 0.0
        val_correct = 0
        val_total = 0
        
        with torch.no_grad():
            for inputs, targets in val_loader:
                inputs = inputs.to(device, non_blocking=True)
                targets = targets.to(device, non_blocking=True)
                
                outputs = model(inputs)
                loss = criterion(outputs, targets)
                
                val_loss += loss.item() * inputs.size(0)
                _, predicted = outputs.max(1)
                val_total += targets.size(0)
                val_correct += predicted.eq(targets).sum().item()
        
        val_loss = val_loss / val_total
        val_acc = 100.0 * val_correct / val_total
        
        # ─── Update LR scheduler ────────────────────────────────
        scheduler.step()
        current_lr = optimizer.param_groups[0]['lr']
        
        # ─── Logging ────────────────────────────────────────────
        epoch_time = time.time() - epoch_start
        
        history['train_loss'].append(train_loss)
        history['val_loss'].append(val_loss)
        history['train_acc'].append(train_acc)
        history['val_acc'].append(val_acc)
        history['lr'].append(current_lr)
        
        # TensorBoard
        writer.add_scalar('Loss/train', train_loss, epoch)
        writer.add_scalar('Loss/val', val_loss, epoch)
        writer.add_scalar('Accuracy/train', train_acc, epoch)
        writer.add_scalar('Accuracy/val', val_acc, epoch)
        writer.add_scalar('LR', current_lr, epoch)
        
        print(f"Epoch {epoch+1:3d}/{n_epochs} [{epoch_time:.1f}s] | "
              f"Train: loss={train_loss:.4f}, acc={train_acc:.1f}% | "
              f"Val: loss={val_loss:.4f}, acc={val_acc:.1f}% | "
              f"LR: {current_lr:.2e}")
        
        # ─── Checkpointing ──────────────────────────────────────
        if val_acc > best_val_acc:
            best_val_acc = val_acc
            best_epoch = epoch + 1
            
            checkpoint = {
                'epoch': epoch + 1,
                'model_state_dict': model.state_dict(),
                'optimizer_state_dict': optimizer.state_dict(),
                'scheduler_state_dict': scheduler.state_dict(),
                'val_acc': val_acc,
                'val_loss': val_loss,
            }
            torch.save(checkpoint, os.path.join(checkpoint_dir, 'best_model.pth'))
            print(f"  *** Saved best model: {val_acc:.2f}% ***")
    
    # ─── Cleanup ──────────────────────────────────────────────────
    writer.close()
    
    print(f"\nTraining complete. Best val acc: {best_val_acc:.2f}% at epoch {best_epoch}")
    
    return best_val_acc, history


def resume_training(model, optimizer, scheduler, checkpoint_path):
    """Resume training from a saved checkpoint."""
    checkpoint = torch.load(checkpoint_path, map_location='cpu')
    
    model.load_state_dict(checkpoint['model_state_dict'])
    optimizer.load_state_dict(checkpoint['optimizer_state_dict'])
    scheduler.load_state_dict(checkpoint['scheduler_state_dict'])
    
    start_epoch = checkpoint['epoch']
    
    print(f"Resumed from epoch {start_epoch}, val_acc={checkpoint['val_acc']:.2f}%")
    return start_epoch
```

---

## 8. GPU Utilization

### Basic GPU Operations

```python
import torch

# Check GPU availability
print(torch.cuda.is_available())           # True if NVIDIA GPU + CUDA installed
print(torch.cuda.device_count())           # number of GPUs
print(torch.cuda.get_device_name(0))       # "NVIDIA RTX 3090" etc.

# M1/M2 Mac:
print(torch.backends.mps.is_available())   # Metal Performance Shaders

# Select device
device = torch.device('cuda' if torch.cuda.is_available()
                       else 'mps' if torch.backends.mps.is_available()
                       else 'cpu')
print(f"Using: {device}")

# Move tensors and models to device
x = torch.randn(100, 784)
x = x.to(device)        # move to GPU
x = x.cuda()            # equivalent for CUDA only

model = SomeModel()
model = model.to(device)   # move all parameters to GPU

# Move back to CPU (for numpy conversion, plotting):
x_cpu = x.cpu()
x_np = x_cpu.detach().numpy()

# Check where a tensor lives:
print(x.device)          # cuda:0  or  cpu
```

### Memory Management

```python
# Common CUDA out of memory errors:
# RuntimeError: CUDA out of memory. Tried to allocate X GB...

# Solutions:
# 1. Reduce batch size
# 2. Use gradient accumulation (simulate larger batch)
# 3. Use mixed precision (FP16 — halves memory)
# 4. Use gradient checkpointing (trade compute for memory)

# Check memory usage:
print(torch.cuda.memory_allocated() / 1024**3, "GB allocated")
print(torch.cuda.max_memory_allocated() / 1024**3, "GB peak")

# Free memory (useful in interactive sessions):
del unused_tensor
torch.cuda.empty_cache()   # releases cached memory back to OS

# Gradient checkpointing (recompute activations during backward instead of storing)
from torch.utils.checkpoint import checkpoint

class MemoryEfficientBlock(nn.Module):
    def __init__(self, ...):
        ...
    
    def forward(self, x):
        # Use gradient checkpointing for large blocks
        return checkpoint(self._forward, x)
    
    def _forward(self, x):
        # actual computation
        ...
```

---

## 9. Model Saving and Loading

```python
import torch

model = SomeModel()
optimizer = torch.optim.Adam(model.parameters())

# ─── Save state dict (PREFERRED) ─────────────────────────────────
# Only saves the parameters, not the model architecture.
# You need to re-create the model object before loading.

torch.save(model.state_dict(), 'model_weights.pth')

# Load:
model_new = SomeModel()   # must re-create architecture first
model_new.load_state_dict(torch.load('model_weights.pth', map_location='cpu'))
model_new.eval()

# ─── Save entire model (NOT recommended for production) ───────────
# Saves architecture + weights using pickle.
# Fragile: breaks if you rename or move the model class.

torch.save(model, 'full_model.pth')
loaded_model = torch.load('full_model.pth', map_location='cpu')

# ─── Save training checkpoint (best practice) ─────────────────────
checkpoint = {
    'epoch': 15,
    'model_state_dict': model.state_dict(),
    'optimizer_state_dict': optimizer.state_dict(),
    'train_loss': 0.234,
    'val_accuracy': 0.923,
    'model_config': {'hidden_dim': 512, 'num_classes': 10}
}
torch.save(checkpoint, 'checkpoint_epoch15.pth')

# Load checkpoint and resume:
checkpoint = torch.load('checkpoint_epoch15.pth', map_location='cpu')
model.load_state_dict(checkpoint['model_state_dict'])
optimizer.load_state_dict(checkpoint['optimizer_state_dict'])
start_epoch = checkpoint['epoch']
print(f"Resuming from epoch {start_epoch}")

# ─── Why state_dict is preferred ─────────────────────────────────
# - model.state_dict() is just a dict of {parameter_name: tensor}
# - Compatible across PyTorch versions
# - Portable: load on different device (map_location)
# - Safe: no arbitrary code execution (unlike pickle of full model)
```

### torchinfo — Model Summary

```python
# pip install torchinfo
from torchinfo import summary

model = SomeModel()

# Print layer-by-layer summary with parameter counts and shapes
summary(
    model,
    input_size=(1, 3, 224, 224),  # single sample shape
    col_names=["input_size", "output_size", "num_params", "trainable"],
    row_settings=["var_names"],
    depth=5
)

# Output:
# ─────────────────────────────────────────────────────────────────────
# Layer (type:depth-idx)    Input Shape      Output Shape    Param #
# ─────────────────────────────────────────────────────────────────────
# ResNet                    [1, 3, 224, 224] [1, 1000]       --
# ├─ Conv2d: 1-1            [1, 3, 224, 224] [1, 64, 112...]  9,408
# ...
# Total params: 25,557,032
# Trainable params: 25,557,032
# ─────────────────────────────────────────────────────────────────────
```

---

## 10. Mixed Precision Training

Modern GPUs (Volta and later from NVIDIA, all Apple Silicon) have hardware support for FP16 or BF16 arithmetic, which is 2-4x faster than FP32 and uses half the memory.

### FP32 vs FP16 vs BF16

```
FP32 (float32): sign(1) + exponent(8) + mantissa(23) = 32 bits
  Range: ±3.4×10^38
  Precision: ~7 decimal digits

FP16 (float16): sign(1) + exponent(5) + mantissa(10) = 16 bits
  Range: ±65504     ← small range! overflow/underflow risk
  Precision: ~3 decimal digits

BF16 (bfloat16): sign(1) + exponent(8) + mantissa(7) = 16 bits
  Range: same as FP32  ← same exponent!
  Precision: ~2 decimal digits
  
BF16 is generally preferred over FP16 for training:
  - Same dynamic range as FP32 → no overflow
  - Less precision loss issues in practice
  - Native on A100, H100, Apple M1/M2
```

### Using torch.cuda.amp

```python
import torch
import torch.cuda.amp as amp

# GradScaler: prevents FP16 underflow by scaling the loss
scaler = amp.GradScaler()

for inputs, targets in train_loader:
    inputs = inputs.to('cuda')
    targets = targets.to('cuda')
    
    optimizer.zero_grad()
    
    # autocast: automatically uses FP16/BF16 for eligible operations
    with amp.autocast(device_type='cuda'):   # or 'cpu' on some systems
        outputs = model(inputs)              # forward in FP16
        loss = criterion(outputs, targets)   # loss computation in FP32
    
    # Scale loss to prevent FP16 underflow
    scaler.scale(loss).backward()
    
    # Unscale gradients before clipping (important!)
    scaler.unscale_(optimizer)
    torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
    
    # Update weights (unscales internally, then calls optimizer.step())
    scaler.step(optimizer)
    
    # Update the scale factor for next iteration
    scaler.update()

# Benefits: ~2x speedup, ~2x memory reduction
# Cost: slight risk of numerical issues (mitigated by GradScaler)
```

---

## 11. Debugging PyTorch

### Common Errors and Fixes

```python
# ─── Shape errors (most common) ───────────────────────────────────

# Error: mat1 and mat2 shapes cannot be multiplied (64×512 and 256×10)
# Cause: wrong input/output dimensions in Linear layer
# Fix: print shapes at each step

class DebugModel(nn.Module):
    def forward(self, x):
        print(f"Input: {x.shape}")        # add these temporarily
        x = self.layer1(x)
        print(f"After layer1: {x.shape}")
        x = self.layer2(x)
        print(f"After layer2: {x.shape}")
        return x

# ─── Loss is NaN ─────────────────────────────────────────────────
# Causes:
# 1. Learning rate too large → gradients explode → NaN
# 2. log(0) in loss computation
# 3. Division by zero (e.g., batch norm with zero variance)

# Debug: add gradient checks
for name, param in model.named_parameters():
    if param.grad is not None:
        if torch.isnan(param.grad).any():
            print(f"NaN gradient in {name}")
        if torch.isinf(param.grad).any():
            print(f"Inf gradient in {name}")

# ─── Forgot model.eval() ─────────────────────────────────────────
# Symptom: validation loss/accuracy is inconsistent or too low
# Fix: always call model.eval() before validation loop

# ─── Forgot optimizer.zero_grad() ────────────────────────────────
# Symptom: gradients accumulate → parameters update too aggressively
# Fix: call optimizer.zero_grad() at the START of each training step

# ─── Wrong loss function ─────────────────────────────────────────
# CrossEntropyLoss expects raw logits (NO softmax applied before)
# BCELoss expects probabilities (AFTER sigmoid)
# BCEWithLogitsLoss expects raw logits (includes sigmoid internally — preferred)

# ─── Device mismatch ─────────────────────────────────────────────
# Error: Expected all tensors to be on the same device
# Fix: ensure inputs AND targets are on the same device as model

device = next(model.parameters()).device   # get model's device
inputs = inputs.to(device)
targets = targets.to(device)

# ─── Hooks for debugging activation statistics ───────────────────
activation_stats = {}

def hook_fn(module, input, output):
    activation_stats[module.__class__.__name__] = {
        'mean': output.mean().item(),
        'std': output.std().item(),
        'has_nan': torch.isnan(output).any().item()
    }

# Register hook on a specific layer:
hook = model.layer1.register_forward_hook(hook_fn)

# Run forward pass
output = model(sample_input)

# Check stats
print(activation_stats)

# Remove hook when done
hook.remove()
```

---

## 12. Summary

```
PYTORCH ESSENTIALS
│
├── Tensors: n-dimensional arrays on CPU or GPU
│     torch.tensor(), torch.randn(), torch.zeros()
│     Operations: +, -, *, /, @, .sum(), .max(), .reshape()
│     Shared memory with NumPy via .from_numpy() / .numpy()
│
├── Autograd: automatic gradient computation
│     requires_grad=True → tracked for differentiation
│     .backward() → computes all gradients
│     .grad → gradient stored here
│     torch.no_grad() → skip gradient tracking (inference)
│     ALWAYS call optimizer.zero_grad() before .backward()
│
├── nn.Module: base class for all networks
│     __init__: define layers
│     forward: define computation
│     model.train() / model.eval() → CRITICAL for Dropout & BatchNorm
│
├── Built-in layers
│     nn.Linear, nn.Conv2d, nn.ReLU, nn.Dropout, nn.BatchNorm2d
│     nn.Embedding (text), nn.Sequential (quick stacking)
│
├── Data pipeline
│     Dataset: implement __len__ and __getitem__
│     DataLoader: batching, shuffling, parallel loading
│     Transforms: Compose, ToTensor, Normalize, RandomFlip
│
├── Training loop pattern:
│     model.train()
│     optimizer.zero_grad()
│     output = model(inputs)
│     loss = criterion(output, targets)
│     loss.backward()
│     clip_grad_norm_(...)
│     optimizer.step()
│
└── Utilities
      torch.save / torch.load (state_dict preferred)
      torch.cuda.amp (mixed precision)
      TensorBoard (SummaryWriter)
      torchinfo (model summary)
```

### Key API Reference

| Operation | Code |
|-----------|------|
| Create tensor | `torch.randn(3, 4)` |
| Move to GPU | `tensor.to('cuda')` |
| Matrix multiply | `a @ b` or `torch.mm(a, b)` |
| Disable grad | `with torch.no_grad():` |
| Backward pass | `loss.backward()` |
| Zero gradients | `optimizer.zero_grad()` |
| Update weights | `optimizer.step()` |
| Save model | `torch.save(model.state_dict(), path)` |
| Load model | `model.load_state_dict(torch.load(path))` |
| Train mode | `model.train()` |
| Eval mode | `model.eval()` |

---

## Mini Projects

### Mini Project 1: Autograd Mechanics Explorer

Understand how PyTorch's autograd works by building computation graphs and inspecting gradients.

**Objective:** Demystify the `backward()` call — see the gradient flow through every operation.

```python
import torch
import torch.nn as nn
import matplotlib.pyplot as plt
import numpy as np

# Part 1: Manual gradient verification
print("=== Part 1: Manual Gradient Verification ===")
x = torch.tensor([2.0, 3.0], requires_grad=True)
y = torch.tensor([4.0, 5.0], requires_grad=True)

# z = sum(x^2 * y + sin(y))
z = (x**2 * y + torch.sin(y)).sum()
z.backward()

# Analytical: dz/dx = 2x*y, dz/dy = x^2 + cos(y)
dz_dx_analytical = 2 * x.detach() * y.detach()
dz_dy_analytical = x.detach()**2 + torch.cos(y.detach())

print(f"  dz/dx computed:   {x.grad.numpy()}")
print(f"  dz/dx analytical: {dz_dx_analytical.numpy()}")
print(f"  dz/dy computed:   {y.grad.numpy()}")
print(f"  dz/dy analytical: {dz_dy_analytical.numpy()}")
print(f"  Gradients match: {torch.allclose(x.grad, dz_dx_analytical)}")

# Part 2: Gradient flow through a custom function
print("\n=== Part 2: Gradient Flow Visualization ===")

class GradientMonitor:
    def __init__(self, name):
        self.name = name
        self.grad_norms = []

    def hook(self, grad):
        self.grad_norms.append(grad.norm().item())
        return grad  # must return grad (or None to keep unchanged)

# Build a deep network and monitor gradient norms per layer
torch.manual_seed(42)
n_layers = 8
monitors = []

def build_deep_net(activation):
    layers = []
    monitors = []
    x = torch.randn(32, 10, requires_grad=False)
    a = x.detach().requires_grad_(True)
    layer_inputs = [a]
    weights = []

    for i in range(n_layers):
        w = torch.randn(10, 10, requires_grad=True) * 0.5
        b = torch.zeros(10, requires_grad=True)
        weights.append((w, b))
        z = layer_inputs[-1] @ w + b
        a = activation(z)
        m = GradientMonitor(f"Layer {i+1}")
        monitors.append(m)
        a.register_hook(m.hook)
        layer_inputs.append(a)

    loss = layer_inputs[-1].sum()
    loss.backward()
    return [m.grad_norms[0] if m.grad_norms else 0.0 for m in monitors]

sigmoid_grads = build_deep_net(torch.sigmoid)
relu_grads    = build_deep_net(torch.relu)
tanh_grads    = build_deep_net(torch.tanh)

fig, axes = plt.subplots(1, 2, figsize=(13, 5))

x_layers = range(1, n_layers+1)
axes[0].plot(x_layers, sigmoid_grads[::-1], 'r-o', label='Sigmoid', markersize=5)
axes[0].plot(x_layers, relu_grads[::-1],    'b-o', label='ReLU',    markersize=5)
axes[0].plot(x_layers, tanh_grads[::-1],    'g-o', label='Tanh',    markersize=5)
axes[0].set_xlabel("Layer (output → input)")
axes[0].set_ylabel("Gradient Norm")
axes[0].set_title("Vanishing Gradient: Gradient Norm vs Layer\n(watch how sigmoid gradients shrink!)")
axes[0].legend()
axes[0].grid(True, alpha=0.3)
axes[0].set_yscale('log')

# Part 3: requires_grad and no_grad effects on memory
print("\n=== Part 3: Memory and requires_grad ===")
torch.manual_seed(0)
model = nn.Sequential(nn.Linear(100, 256), nn.ReLU(), nn.Linear(256, 10))

# Inference without grad (no gradient tape)
x_input = torch.randn(64, 100)
with torch.no_grad():
    out_no_grad = model(x_input)
print(f"  no_grad output requires_grad: {out_no_grad.requires_grad}")

# With grad (training mode)
out_with_grad = model(x_input)
print(f"  with_grad output requires_grad: {out_with_grad.requires_grad}")

# Gradient accumulation example
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)
losses = []
for step in range(50):
    out = model(x_input)
    loss = out.pow(2).mean()
    loss.backward()
    # Accumulate gradients for 5 steps then update (simulates larger batch)
    if (step + 1) % 5 == 0:
        optimizer.step()
        optimizer.zero_grad()
    losses.append(loss.item())

axes[1].plot(losses, 'purple', linewidth=1.5)
for i in range(4, 50, 5):
    axes[1].axvline(i, color='red', alpha=0.3, linestyle='--')
axes[1].set_xlabel("Step")
axes[1].set_ylabel("Loss")
axes[1].set_title("Gradient Accumulation\n(red lines = optimizer steps every 5 iterations)")
axes[1].grid(True, alpha=0.3)

plt.tight_layout()
plt.savefig("autograd_mechanics.png", dpi=150)
plt.show()
print("Saved: autograd_mechanics.png")
```

---

### Mini Project 2: Custom Dataset + DataLoader Pipeline

Build a complete data loading pipeline with custom transforms, caching, and augmentation.

**Objective:** Learn the PyTorch data loading pattern — the foundation of every real training loop.

```python
import torch
from torch.utils.data import Dataset, DataLoader
import numpy as np
import matplotlib.pyplot as plt
from sklearn.datasets import make_classification

class SyntheticTabularDataset(Dataset):
    """Tabular classification dataset with optional augmentation."""
    def __init__(self, X, y, augment=False, noise_std=0.05):
        self.X = torch.FloatTensor(X)
        self.y = torch.LongTensor(y)
        self.augment  = augment
        self.noise_std = noise_std

    def __len__(self):
        return len(self.y)

    def __getitem__(self, idx):
        x = self.X[idx].clone()
        if self.augment:
            x += torch.randn_like(x) * self.noise_std  # Gaussian noise augmentation
        return x, self.y[idx]

    @classmethod
    def from_sklearn(cls, **kwargs):
        X, y = make_classification(**kwargs)
        # Standardize
        X = (X - X.mean(0)) / (X.std(0) + 1e-8)
        return cls(X, y)

# Generate data
np.random.seed(42)
dataset_full = SyntheticTabularDataset.from_sklearn(
    n_samples=2000, n_features=20, n_informative=10, n_redundant=5,
    n_classes=3, random_state=42
)

# Split into train/val/test
n = len(dataset_full)
n_train = int(0.7 * n); n_val = int(0.15 * n)
train_ds, val_ds, test_ds = torch.utils.data.random_split(
    dataset_full, [n_train, n_val, n - n_train - n_val],
    generator=torch.Generator().manual_seed(42)
)

# Override augment only for train
train_ds.dataset.augment = True

train_loader = DataLoader(train_ds, batch_size=64, shuffle=True,  num_workers=0, pin_memory=False)
val_loader   = DataLoader(val_ds,   batch_size=256, shuffle=False, num_workers=0)
test_loader  = DataLoader(test_ds,  batch_size=256, shuffle=False, num_workers=0)

print(f"Dataset sizes: train={len(train_ds)}, val={len(val_ds)}, test={len(test_ds)}")
print(f"Train batches: {len(train_loader)}, Val batches: {len(val_loader)}")

# Verify a batch
X_batch, y_batch = next(iter(train_loader))
print(f"Batch shape: X={X_batch.shape}, y={y_batch.shape}")
print(f"Label distribution in batch: {torch.bincount(y_batch).tolist()}")

# Quick training loop to test the pipeline
class TabularNet(nn.Module):
    def __init__(self, n_in, n_classes):
        super().__init__()
        self.net = nn.Sequential(
            nn.Linear(n_in, 128), nn.BatchNorm1d(128), nn.ReLU(), nn.Dropout(0.3),
            nn.Linear(128, 64),  nn.BatchNorm1d(64),  nn.ReLU(), nn.Dropout(0.2),
            nn.Linear(64, n_classes)
        )
    def forward(self, x): return self.net(x)

import torch.nn as nn
torch.manual_seed(42)
model = TabularNet(20, 3)
optimizer = torch.optim.Adam(model.parameters(), lr=0.001)
criterion = nn.CrossEntropyLoss()

train_losses, val_losses, val_accs = [], [], []
for epoch in range(30):
    model.train()
    epoch_loss = 0
    for X_b, y_b in train_loader:
        optimizer.zero_grad()
        loss = criterion(model(X_b), y_b)
        loss.backward()
        optimizer.step()
        epoch_loss += loss.item()
    train_losses.append(epoch_loss / len(train_loader))

    model.eval()
    with torch.no_grad():
        val_loss = sum(criterion(model(X_b), y_b).item() for X_b, y_b in val_loader) / len(val_loader)
        correct = sum((model(X_b).argmax(1) == y_b).sum().item() for X_b, y_b in val_loader)
        val_acc = correct / len(val_ds)
    val_losses.append(val_loss)
    val_accs.append(val_acc)

fig, axes = plt.subplots(1, 2, figsize=(12, 4))
axes[0].plot(train_losses, label='Train'); axes[0].plot(val_losses, label='Val')
axes[0].set_title("Loss Curves"); axes[0].set_xlabel("Epoch"); axes[0].legend(); axes[0].grid(True, alpha=0.3)
axes[1].plot(val_accs, 'green'); axes[1].set_title("Validation Accuracy")
axes[1].set_xlabel("Epoch"); axes[1].grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig("pytorch_pipeline.png", dpi=150)
plt.show()
print(f"Final test accuracy: {val_acc:.3f}")
```

---

## Exercises

1. **Tensor operations**: Create two (4×3) tensors. Compute their matrix product, element-wise product, sum of each row, and L2 norm. Verify shapes at each step.

2. **Autograd from scratch**: Compute the gradient of `f(x,y) = x²y + e^(xy)` at `(x=1, y=2)` using PyTorch autograd. Verify with analytical computation.

3. **Custom layer**: Implement `nn.Linear` from scratch as a custom `nn.Module` using `nn.Parameter`. Verify that it gives the same output as the built-in one with the same weights.

4. **Dataset pipeline**: Create a custom Dataset that loads the MNIST CSVs. Implement `__len__` and `__getitem__`. Build a DataLoader and verify batch shapes.

5. **Training loop bugs**: Take the MNIST training loop. Intentionally introduce 3 bugs: remove `optimizer.zero_grad()`, remove `model.eval()` during validation, and apply softmax before `CrossEntropyLoss`. Run with each bug and observe what goes wrong. Fix each.

6. **Mixed precision**: Enable AMP on the MNIST MLP. Measure training time and GPU memory with and without AMP. Report speedup.

---

**Chapter Summary:**

PyTorch's core innovation is the dynamic computation graph — the graph is built on-the-fly as operations execute, making it feel like regular Python. Tensors are the fundamental unit, supporting all array operations on both CPU and GPU with seamless conversion to/from NumPy. Autograd tracks operations on tensors with `requires_grad=True` and computes gradients via `.backward()`. The `nn.Module` class provides structure for building networks: `__init__` defines components, `forward` defines computation. The training loop follows a strict pattern (zero_grad → forward → loss → backward → step), and model.train() / model.eval() must be called appropriately. Mixed precision training (torch.cuda.amp) gives 2x speedup at minimal accuracy cost, and proper checkpointing (saving state_dict) enables reliable training resumption.

---

**What's Next →** [Chapter 18: Convolutional Neural Networks](./18-convolutional-neural-networks.md)

*Fully connected networks treat every pixel as independent. Convolutional networks understand that adjacent pixels are related — and this inductive bias is why they achieve superhuman performance on images with far fewer parameters.*
