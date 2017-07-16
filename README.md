### plain.im

plain.im is a plaintext paste service with 24 hours expiration time.

#### API

Paste the command output:

```command | curl -sF "plain=<-" https://plain.im/```

Create an alias:

```alias plain='curl -sF "plain=<-" https://plain.im/'```

to shorten paste command to just

```command | plain```

Paste and make link available as an X selection for pasting into X applications (requires xclip):

```alias plain='curl -sF "plain=<-" https://plain.im/ | tee /dev/stderr | xclip'```

#### PASTE REMOVAL

Paste can be removed by making a DELETE request to the paste URL:

```curl -X DELETE https://plain.im/key'```

#### LICENSE

MIT
