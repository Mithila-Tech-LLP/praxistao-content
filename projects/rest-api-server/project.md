---
title: Build a REST API Server
subtitle: Learn HTTP, routing, middleware, and CRUD by building a real REST API
category: Web APIs
difficulty: intermediate
duration: 4-6 hours
accent: "#fbbf24"
technologies: [Go]
skills: [HTTP, Routing, Middleware, JSON, REST, In-Memory Store]
prerequisites: [basic-programming]
repo: rest-api-server
outcomes:
  - Build an HTTP server with multiple routes
  - Parse URL path parameters and query strings
  - Read and write JSON request and response bodies
  - Implement CRUD with an in-memory store
  - Write reusable middleware for logging and auth
  - Handle HTTP errors with proper status codes
---

## Overview

REST (Representational State Transfer) is the architectural style behind nearly every web API you interact with — GitHub, Stripe, Twitter, and thousands more. In this project you will build one from scratch using only Go's standard library.

You will start with a single handler that returns JSON, and by the end you will have a fully working REST API with routing, middleware, an in-memory data store, and proper error handling — all without a single external dependency.

## Who Should Take This

You have completed Basic Programming (or equivalent) and want to understand how web APIs work under the hood. You are comfortable with functions, structs, and maps in Go.

## What You Will Build

A REST API for managing a simple "todos" resource: create, read, update, delete todo items stored in memory.
