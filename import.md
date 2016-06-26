# Data import things

### Parse Station lines from some Wikipedia

    ^([A-Z]+)-?(\d\d)\s+([-\w]+)\s+([^-\s\x3400-\x4DB5\x4E00-\x9FCB\xF900-\xFA6A]+)\s+(\d\.\d).*$

#### To nodes

    \1\2 [label="{{\4|\3}|\1\2}"];

#### To edges

    \1\2 -> \1\2 [len=\5];
