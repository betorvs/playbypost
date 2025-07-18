# Gemini Development Workflow

This document outlines the standard process for using Gemini to manage software development projects, based on the rules and workflows defined in the `.cursor/rules` directory.

## Core Principles

- **Task-Driven Development**: The workflow is centered around a `tasks.json` file, which is managed by a `task-master` tool.
- **Structured Rules**: All development activities should adhere to the guidelines specified in the `.cursor/rules` directory.
- **Continuous Improvement**: The rules and workflows are designed to be improved over time based on new patterns and best practices.

## The Basic Loop

The fundamental development cycle is as follows:

1.  **`list`**: Show what needs to be done.
2.  **`next`**: Decide what to work on.
3.  **`show <id>`**: Provide details for a specific task.
4.  **`expand <id>`**: Break down a complex task into smaller, manageable subtasks.
5.  **Implement**: Write the code and tests.
6.  **`update-subtask`**: Log progress and findings.
7.  **`set-status`**: Mark tasks and subtasks as `done`.
8.  **Repeat**.

## Taskmaster Tool

The `task-master` tool is the primary interface for managing tasks. It can be used via the command line or through an MCP server for integrated environments.

### Key Commands

-   **`task-master init`**: Initialize a new project.
-   **`task-master parse-prd [file]`**: Parse a Product Requirements Document (PRD) to generate initial tasks.
-   **`task-master list`**: List all tasks.
-   **`task-master next`**: Show the next available task.
-   **`task-master show [id]`**: Show details for a specific task.
-   **`task-master add-task`**: Add a new task.
-   **`task-master expand --id=<id>`**: Break down a task into subtasks.
-   **`task-master set-status --id=<id> --status=<status>`**: Update the status of a task.

### Tagged Task Lists

Taskmaster supports tagged task lists to manage different contexts (e.g., features, branches, experiments).

-   **`task-master tags`**: List all available tags.
-   **`task-master add-tag <tagName>`**: Create a new tag.
-   **`task-master use-tag <tagName>`**: Switch to a different tag.

## Rule-Based Guidance

Gemini's behavior is guided by a set of rules defined in `.mdc` files. These rules cover:

-   **Code Style and Patterns**: Enforced through `cursor_rules.mdc`.
-   **Self-Improvement**: Guidelines for updating and improving the rules themselves, as detailed in `self_improve.mdc`.
-   **Development Workflow**: The overall development process is described in `dev_workflow.mdc`.
-   **Taskmaster Commands**: A comprehensive reference for all `task-master` commands is available in `taskmaster.mdc`.

## Getting Started

1.  Initialize the project with `task-master init`.
2.  Create a PRD and parse it with `task-master parse-prd` to generate the initial set of tasks.
3.  Use `task-master list` and `task-master next` to start working on tasks.
4.  Follow the iterative development loop, using the various `task-master` commands to manage your workflow.
