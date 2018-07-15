echo 'strict digraph {'
echo '  edge [style=bold color="#000000" fontname=Arial]'
echo '  graph [sep=10 bgcolor="#666666"]'
echo '  node [shape=record style="filled,bold" fontname=Arial fillcolor="#ffffff"]'
cat line-*.dot
echo '  edge [style=dashed color="#dddddd"]'
echo '  node [shape=diamond color="#ff0000"]'
cat transfers.dot
echo '}'
