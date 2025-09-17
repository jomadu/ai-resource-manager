# TypeScript Coding Guidelines

Best practices for TypeScript development in our projects.

## Type Safety Rules

- Use strict type checking in tsconfig.json
- Prefer interfaces over types for object shapes
- Avoid any type, use unknown instead
- Always handle promise rejections

```typescript
// Good example
interface User {
  name: string;
  email: string;
}

function getUser(): Promise<User> {
  return fetch('/api/user').then(res => res.json());
}
```

```typescript
// Bad example
function getUser(): any {
  return fetch('/api/user').then(res => res.json());
}
```

## Code Style

- Use 2 spaces for indentation
- Use single quotes for strings
- Add trailing commas in multiline objects

```typescript
const config = {
  name: 'myApp',
  version: '1.0.0',
};
```
