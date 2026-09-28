# Subtitles_Visualizer
## Summary
This project is supposed to provide a minimal overlay for rendering subtitles generated in realtime and fed to it through stdin.
Usage is meant to be  
``subtitle_generator | subtitle_visualizer``
with subtitle_generator being any suitable such program.

## Installation
The project only depends on [go-gui](https://www.github.com/go-gui-org/go-gui),
which means you simply need to type
``go build .``
and the entire dependency tree will be downloaded and built together with the program.
