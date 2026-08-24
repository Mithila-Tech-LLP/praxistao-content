# Chapter 22: TypeScript Advanced — Generics, Utility Types & Conditional Types

TypeScript is the typed superset of JavaScript used in most serious Node.js projects. Senior engineers are expected to know advanced TypeScript patterns, not just basic type annotations.

## Table of Contents

1. [Generics with Constraints](#1-generics-with-constraints)
2. [Utility Types — The 8 You Must Know](#2-utility-types--the-8-you-must-know)
3. [Mapped Types](#3-mapped-types)
4. [Conditional Types](#4-conditional-types)
5. [Template Literal Types](#5-template-literal-types)
6. [Discriminated Unions](#6-discriminated-unions)
7. [Summary](#summary)

---

## 1. Generics with Constraints

```typescript
// Basic generic
function identity<T>(arg: T): T { return arg; }

// Generic with constraint: T must have a .length property
function logLength<T extends { length: number }>(arg: T): T {
    console.log(arg.length);
    return arg;
}
logLength("hello"); // works: string has length
logLength([1,2,3]); // works: array has length
logLength(42);      // ERROR: number has no length

// Constrain to keys of an object (keyof)
function getProperty<T, K extends keyof T>(obj: T, key: K): T[K] {
    return obj[key];
}
const user = { name: "Alice", age: 30 };
const name = getProperty(user, "name"); // type: string
const age  = getProperty(user, "age");  // type: number
getProperty(user, "email");             // ERROR: not a key
```

---

## 2. Utility Types — The 8 You Must Know

```typescript
interface User {
    id: string;
    name: string;
    email: string;
    age: number;
}

// Partial<T>: makes all properties optional
type PartialUser = Partial<User>;
// { id?: string; name?: string; email?: string; age?: number }
function updateUser(id: string, updates: Partial<User>) {}

// Required<T>: makes all properties required (removes optional)
type RequiredUser = Required<PartialUser>;

// Readonly<T>: makes all properties readonly
type ReadonlyUser = Readonly<User>;

// Pick<T, K>: pick subset of properties
type UserPreview = Pick<User, "id" | "name">;
// { id: string; name: string }

// Omit<T, K>: omit subset of properties
type UserWithoutAge = Omit<User, "age">;
// { id: string; name: string; email: string }

// Record<K, V>: object type with keys K and values V
type UserMap = Record<string, User>;
const users: UserMap = { "alice": { id: "1", name: "Alice", email: "a@b.com", age: 30 } };

// Exclude<T, U>: exclude types from a union
type A = "a" | "b" | "c";
type BC = Exclude<A, "a">; // "b" | "c"

// Extract<T, U>: extract types from a union that are assignable to U
type AOrB = Extract<A, "a" | "b">; // "a" | "b"

// NonNullable<T>: removes null and undefined
type MaybeString = string | null | undefined;
type DefiniteString = NonNullable<MaybeString>; // string

// ReturnType<T>: get return type of a function
function getUser() { return { name: "Alice", age: 30 }; }
type UserType = ReturnType<typeof getUser>; // { name: string; age: number }

// Parameters<T>: get parameter types of a function
type GetUserParams = Parameters<typeof getUser>; // []
```

---

## 3. Mapped Types

Mapped types create new types by transforming each property of an existing type.

```typescript
// Make all properties nullable
type Nullable<T> = { [K in keyof T]: T[K] | null };
type NullableUser = Nullable<User>;
// { id: string | null; name: string | null; ... }

// Add 'get' prefix to all properties (common pattern for validation schemas)
type Getters<T> = { [K in keyof T as `get${Capitalize<string & K>}`]: () => T[K] };
type UserGetters = Getters<User>;
// { getId: () => string; getName: () => string; ... }

// Filter properties by type
type OnlyStrings<T> = {
    [K in keyof T as T[K] extends string ? K : never]: T[K]
};
type UserStringFields = OnlyStrings<User>;
// { id: string; name: string; email: string } — age (number) is excluded
```

---

## 4. Conditional Types

```typescript
// Basic conditional type: T extends U ? Yes : No
type IsString<T> = T extends string ? true : false;
type A = IsString<string>; // true
type B = IsString<number>; // false

// The 'infer' keyword: extract a type within a conditional
type UnpackPromise<T> = T extends Promise<infer U> ? U : T;
type Resolved = UnpackPromise<Promise<string>>; // string
type Direct = UnpackPromise<number>;            // number

// Extract return type using infer (reimplementation of ReturnType)
type MyReturnType<T extends (...args: any) => any> = 
    T extends (...args: any) => infer R ? R : never;

// Distributive conditional types: distributes over unions
type ToArray<T> = T extends any ? T[] : never;
type StrOrNumArray = ToArray<string | number>;
// string[] | number[]  (distributes over each member of the union)
```

---

## 5. Template Literal Types

```typescript
// Combine string literals
type EventName = "click" | "focus" | "blur";
type HandlerName = `on${Capitalize<EventName>}`;
// "onClick" | "onFocus" | "onBlur"

// Type-safe API routes
type Routes = `/api/users/${string}` | "/api/health";
const route: Routes = "/api/users/123"; // OK
const bad: Routes = "/api/posts";       // ERROR

// Type-safe CSS property builder
type CSSProperty = "margin" | "padding" | "border";
type Side = "top" | "right" | "bottom" | "left";
type CSSBoxProperty = `${CSSProperty}-${Side}`;
// "margin-top" | "margin-right" | ... | "border-left"
```

---

## 6. Discriminated Unions

A pattern for modeling state machines and sum types safely.

```typescript
type LoadingState = { status: "loading" };
type SuccessState = { status: "success"; data: User[] };
type ErrorState   = { status: "error"; error: Error };

type RequestState = LoadingState | SuccessState | ErrorState;

function renderState(state: RequestState) {
    switch (state.status) {  // 'status' is the discriminant
        case "loading":
            return "Loading...";
        case "success":
            return `${state.data.length} users`; // TypeScript knows state.data exists here
        case "error":
            return `Error: ${state.error.message}`; // TypeScript knows state.error exists here
    }
    // TypeScript can verify exhaustiveness here
}

// Exhaustiveness check: TypeScript errors if you miss a case
function assertNever(x: never): never {
    throw new Error("Unexpected value: " + x);
}
// Add to the end of switch: default: return assertNever(state);
```

---

## Summary

- **Generics with constraints:** `T extends SomeInterface` to restrict what types can be used.
- **Utility types:** Partial, Required, Readonly, Pick, Omit, Record, Exclude, Extract, NonNullable, ReturnType, Parameters.
- **Mapped types:** transform each property of a type using `[K in keyof T]`.
- **Conditional types:** `T extends U ? Yes : No`. Use `infer` to extract types.
- **Template literal types:** build string literal types from combinations.
- **Discriminated unions:** model state machines with a discriminant property. Use `switch` for exhaustiveness.
