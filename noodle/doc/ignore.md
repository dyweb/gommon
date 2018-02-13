# Handle ignore spec 

- Made a fork of go.rice w/ ignore support when building Ayi 
  - https://github.com/at15/go.rice/issues/1
  - was using pattern in https://github.com/codeskyblue/dockerignore
  - the official dockerignore implementation https://github.com/moby/moby/blob/master/builder/dockerignore/dockerignore.go
- **only support single ignore file**
- can NOT reuse `.gitignore`
  - the gitignore semantic is more complex than simple wildcards https://git-scm.com/docs/gitignore
- [ ] need to figure out how to handle `*` matching correctly ...
  - it is also a leetcode question https://leetcode.com/problems/wildcard-matching/description/
  
````text
# example of .noodleignore
# support comment, also blank line should be ignored

vendor # this ignore the vendor directory, relative to the path of ignore file
node_modules
assets/*.partial.html # ignore assets/a.partial.html etc.
````