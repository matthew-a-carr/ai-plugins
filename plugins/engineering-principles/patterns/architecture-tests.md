# Pattern — Architecture tests (mechanical enforcement)

Constitution P14 in practice. Encode the architecture rules as tests so the build fails when a rule is broken.

## What to enforce

For a clean-architecture codebase (per P2), the minimum:

1. **Layer imports** — `domain/` imports nothing from `application/`, `infrastructure/`, or `delivery/`. `application/` imports only `domain/`. `infrastructure/` may import `domain/` and `application/` but not the other way around.
2. **Composition root** — only one file may construct adapters (e.g. `Drizzle*Repository`, `Sql*Repository`). Other constructions fail the build.
3. **Framework leak** — no framework imports (`@nestjs/*`, `org.springframework.*`, `next/*`) inside `domain/`.
4. **Async in domain** — domain functions are pure (P13). No `async`, no Promise, no IO type in domain signatures.
5. **Result over throw** — domain functions return `Result<T, E>` (or equivalent). Throwing is restricted to adapter-level "exceptional" failures.

## Shape (JVM example using ArchUnit)

```java
@AnalyzeClasses(packages = "com.example.invoicing")
class LayerBoundariesTest {
    @ArchTest
    static final ArchRule domain_has_no_framework_dependencies = noClasses()
        .that().resideInAPackage("..domain..")
        .should().dependOnClassesThat().resideInAPackage("org.springframework..");

    @ArchTest
    static final ArchRule application_does_not_depend_on_delivery = noClasses()
        .that().resideInAPackage("..application..")
        .should().dependOnClassesThat().resideInAPackage("..delivery..");
}
```

## Shape (TypeScript example using import-rule tests)

```typescript
describe("architecture", () => {
  it("domain has no infrastructure imports", () => {
    const domainFiles = glob.sync("src/domain/**/*.ts");
    for (const f of domainFiles) {
      const source = readFileSync(f, "utf8");
      expect(source).not.toMatch(/from ["']\.\.\/infrastructure/);
      expect(source).not.toMatch(/from ["']next\//);
    }
  });
});
```

## Where it lives

- A test file under the project's standard test path (e.g. `src/__tests__/architecture.test.ts`, `src/test/java/.../architecture/LayerBoundariesTest.java`).
- Runs in the standard unit-test phase. No special invocation.
- A failure is a normal red test, not a special tooling alert.

## Why this matters

Documentation tells humans what to do. Tests tell the build what to do. Only the build can be relied on.

## References

- Constitution P2 (clean architecture), P13 (domain modelling), P14 (mechanical enforcement)
- [ArchUnit](https://www.archunit.org/) for JVM
- [dependency-cruiser](https://github.com/sverweij/dependency-cruiser) / [eslint-plugin-boundaries](https://github.com/javierbrea/eslint-plugin-boundaries) for TypeScript
- madetech/clean-architecture — same rules, same enforcement stance
