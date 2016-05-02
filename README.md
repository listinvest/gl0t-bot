# gl0t-bot

 There are a bunch of different ways to use something as a communication pipe with CNC, such as Google Docs, Gmail, Twitter or Evernote. gl0t-bot is a PoC which is an alternative way that used glot.io (Pastebin-like) service as a command proxy server, to pass along data to CNC. Actually you can use any Pastebin-like service. 

It is currently just a PoC, it’s not finished yet and I’m not intended to use in a real world, Just doing for fun.

# How does it work ?:
- Register a throwaway glot account for a command account and report account
- Generate API Token for 2 accounts and set the variable in the source file (APIToken, ReportAPIToken)
- Let’s start with a non CNC approach, on the command account, create new snippet by choose a plaintext and set title as “ALL:123456” (BOT_ID:COMMAND_ID, if the BOT_ID is “ALL”, every bot will take this command)  Now the content, we will use this simple interface to execute shell command.
```
{"cmd": "<ANY_COMMAND>", "args": "<ANY_ARGS>"}
```
eg.
```
{"cmd": "ifconfig", "args": ""}
```
Encode it to base64 and save
- Run the source file
- On the report account, you should see the created report snippet

# Note:
- Step 3,5 are controlled by CNC, it’s your job to fill the gaps.
- Old snippet is still there even the command successfully executed, it should handle by CNC
