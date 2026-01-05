# Codebase Indexing MCP Setup

## Purpose

To prevent code duplication and maintain codebase awareness, we should use a codebase indexing MCP server.

## Options

### 1. Cursor Built-in Indexing
- **Status**: Cursor IDE has built-in codebase indexing
- **Usage**: Automatic - Cursor indexes the codebase for context
- **Note**: This is already active in Cursor

### 2. MCP Codebase Indexing Servers
- Search for available MCP servers that provide codebase indexing
- May provide additional features like semantic search

## Recommendation

**Use Cursor's built-in indexing** - it's already active and provides:
- Full codebase context awareness
- Semantic search
- Code navigation
- Duplicate detection

## Best Practices to Maximize Indexing Effectiveness

1. **Clear Package Structure**: Well-organized packages help indexing
2. **Descriptive Names**: Clear function/variable names improve search
3. **Comments**: Document public APIs for better context
4. **Consistent Patterns**: Reusable patterns are easier to find

## Action Items

- [x] Cursor indexing is active (built-in)
- [ ] Ensure code follows clear patterns for better indexing
- [ ] Use shared utilities to avoid duplication
- [ ] Keep functions focused and well-named

