# CHANGELOG

## develop

### New

Features:
- added brace expansion
- added tilde expansion
- added parameter expansion
  - added expand-to-value
  - added expand-with-default-value
  - added expand-assign-default-value
  - added expand-write-error
  - added expand-use-alternate-value
  - added expand-to-substring
  - added expand-to-substring-length
  - added expand-prefix-match-names
  - added expand-parameter-length
  - added expand-no-positional-params
  - added expand-remove-shortest-prefix
  - added expand-remove-longest-prefix
  - added expand-remove-shortest-suffix
  - added expand-remove-longest-suffix
  - added expand-search-replace-all-matches
  - added expand-search-replace-first-match
  - added expand-search-replace-prefix
  - added expand-search-replace-suffix
  - added expand-uppercase-first-char
  - added expand-uppercase-all-chars
  - added expand-lowercase-first-char
  - added expand-lowercase-all-chars

Exported API:
- added `Expand()`
- added `ExpandTilde()`
- added `ExpansionCallbacks`

Errors:
- added `ErrMismatchedBrace`
- added `ErrMismatchedClosingBrace`
