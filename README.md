# Inspector Gopher
<img src="https://raw.githubusercontent.com/gophergala2016/inspector_gopher/master/public/inspector_gufer_ready.png" align="right"/>

### Google Trends for your repo

We live in one the best times in history for Developers, especially if you are into opensource (Go in particular). Even though opensource has never been so popular its still very difficult for someone to get started.

One of the reasons is because repositories are too big and its quite difficult to find on which part of the code the community is working on.

### This is where Inspector Gopher comes in.
We analyze your repo down to the function/struct level and make a visualization based on the hotness of that particular part of code.

### Hotness (Volatility)

Our volatility estimation algorithm was inspired by [TF-IDF](http://www.tfidf.com/). The volatility of a function or a struct is a metric that showcases how much that particular part of code is prone to change.

### Result

By looking exactly at this metric, we can plot a Tree view graph that shows contributors which part of code is having the most relevant changes in the recent history.

![Tree Map](https://raw.githubusercontent.com/gophergala2016/inspector_gopher/master/public/treemap.png)