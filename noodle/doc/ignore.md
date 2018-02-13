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

vendor # ignore any file or directory whose name is vendor

# ignore all the files and directory under test,
# since it also applies to walk, test/sub/example.txt will be ignored as well
# however it is not ignored because match test/* pattern, * does not match separator
# TODO: this is just my assumption ... not tested ...
test/*

# ignore assets/a.partial.html etc.
assets/*.partial.html
````