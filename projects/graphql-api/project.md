---
title: Build a GraphQL API
subtitle: Learn schemas, queries, mutations, and nested resolvers with graphql-go
category: Web APIs
difficulty: intermediate
duration: 4-6 hours
accent: "#e879f9"
technologies: [Go, GraphQL]
skills: [Schema, Queries, Mutations, Resolvers, Arguments, Input Types, Errors]
prerequisites: [rest-api-server]
repo: graphql-api
outcomes:
  - Define a GraphQL schema with types and fields
  - Write query resolvers that return data
  - Implement mutations that modify in-memory state
  - Resolve nested/related types
  - Handle arguments and input types
  - Return proper GraphQL errors
---

## Overview

GraphQL is the query language behind GitHub's API v4, Shopify, and thousands of modern web applications. Unlike REST — where the server decides what data to return — GraphQL lets the client ask for exactly what it needs.

In this project you will build a GraphQL API for a simple blog with authors, posts, and comments — all in memory. Each task adds a new GraphQL concept.

**First-run note:** The first time you run `lncli run`, Go will download graphql-go (~15s). Subsequent runs are instant.
