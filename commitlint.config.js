export default { 
    extends: ['@commitlint/config-conventional'],
    rules: {
        "body-leading-blank": [2, "always"],
        "body-empty": [2, "never"],
        "body-max-line-length": [2, "always", 80],
        "footer-leading-blank": [2, "always"],
        "footer-max-line-length": [2, "always", 80],
        "header-max-length": [2, "always", 80],
        "header-trim": [2, "always"],
        "type-enum": [2, "always", [
            "build",
            "chore",
            "ci",
            "docs",
            "feat",
            "fix",
            "perf",
            "refactor",
            "revert",
            "style",
            "test",
        ]],
        "type-empty": [2, "never"],

    }
};
