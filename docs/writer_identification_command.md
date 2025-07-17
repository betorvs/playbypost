# Writer Identification Command Structure

## Overview
This document defines the command structure for users to identify themselves as writers within the Discord and Slack plugins. The command aims to be intuitive and consistent across platforms, facilitating the association between a user's platform account and a writer's account in the Play-by-Post system.

## Command Syntax

The command will follow this general syntax:

```
/iamwriter <writer_username>
```

- **`/iamwriter`**: This is the command trigger. It should be consistent in both Discord and Slack.
- **`<writer_username>`**: This is a required parameter representing the unique username of the writer account the user wishes to associate with. This username must correspond to an existing writer account in the Play-by-Post system.

## Expected Behavior

Upon invocation of the command, the following actions are expected:

1.  **Parameter Extraction**: The plugin (Discord or Slack) will parse the command to extract the provided `<writer_username>`.
2.  **User Identification**: The plugin will identify the `user_id` of the individual who executed the command (e.g., Discord User ID, Slack User ID).
3.  **Backend Communication**: The plugin will send a request to the backend API (which will be implemented in a later subtask) to create an association between the identified `user_id` and the `writer_id` corresponding to the `<writer_username>`.
4.  **User Feedback**: The plugin should provide immediate feedback to the user in the chat, indicating the success or failure of the association attempt:
    *   **Success**: "You are now associated with writer `<writer_username>`."
    *   **Failure (Writer Not Found)**: "Writer `<writer_username>` not found. Please check the username and try again."
    *   **Failure (Already Associated)**: "You are already associated with writer `<writer_username>`."
    *   **Failure (Other Error)**: "An error occurred while trying to associate you with a writer. Please try again later."

## Example Usage

- **Discord/Slack**: `/iamwriter my_storyteller_name`

## Future Considerations

- Implement a command to disassociate from a writer.
- Implement a command to list current associations.
