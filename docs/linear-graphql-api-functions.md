# Linear GraphQL API Functions

This document provides information about the Linear GraphQL API functions including queries, mutations, and subscriptions.

## Overview

Linear's public API is built using GraphQL. It's the same API used internally for developing Linear applications. The API endpoint is `https://api.linear.app/graphql`.

## Resources for Complete API Reference

Since the Linear GraphQL API is extensive and constantly evolving, the most accurate and up-to-date list of all available functions can be found through these official resources:

### 1. Apollo Studio Explorer (Recommended)
- **URL**: https://studio.apollographql.com/public/Linear-API/schema/reference?variant=current
- **Features**: 
  - Interactive schema browser
  - No login required
  - Searchable interface
  - Real-time query testing
- **Usage**: Click the "Schema" tab to browse all available queries, mutations, and subscriptions

### 2. GitHub Schema File
- **URL**: https://github.com/linear/linear/blob/master/packages/sdk/src/schema.graphql
- **Features**: 
  - Complete GraphQL schema definition
  - Always up-to-date
  - Raw schema format

### 3. GraphQL Introspection
- **Endpoint**: `https://api.linear.app/graphql`
- **Method**: Use GraphQL introspection queries with proper authentication

## Common Operation Categories

Based on the Linear API documentation and available MCP functions, here are the main categories of operations:

### Queries
- **Issues**: Fetch issues, issue details, issue statuses, labels
- **Projects**: List and get project information
- **Teams**: List and get team information
- **Users**: List and get user information
- **Documents**: List and get documents
- **Cycles**: List team cycles
- **Comments**: Fetch issue comments
- **Search**: Search across various entities

### Mutations
- **Issues**: Create, update, archive issues
- **Projects**: Create, update projects
- **Comments**: Create comments on issues
- **Documents**: Create, update documents
- **Attachments**: Upload and manage file attachments
- **Webhooks**: Manage webhook subscriptions

### Subscriptions
- Real-time updates for various entities when changes occur

## Authentication

The Linear API supports two authentication methods:
- **Personal API Keys**: For personal use and development
- **OAuth2**: Recommended for applications used by others

## Naming Conventions (SDK v2)

For mutations in the Linear SDK v2, the naming pattern follows:
- Action verb (create, update, delete, archive) precedes the model name
- Example: `createIssue`, `updateProject`, `archiveIssue`

## Rate Limiting

The Linear API implements rate limiting to ensure fair usage. Check the official documentation for current limits.

## Getting Started

1. **For exploration**: Visit the Apollo Studio Explorer link above
2. **For development**: 
   - Install the Linear SDK: `npm install @linear/sdk`
   - Use your API key for authentication
   - Reference the schema for available operations

## Additional Resources

- **Official Documentation**: https://developers.linear.app/
- **Getting Started Guide**: https://linear.app/developers/graphql
- **API and Webhooks**: https://linear.app/docs/api-and-webhooks

## Note

For the most comprehensive and current list of all GraphQL functions, please refer to the Apollo Studio Explorer or the GitHub schema file linked above. The Linear API is actively developed, and new functions are added regularly.