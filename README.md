# ccache-go
A golang wrapper for ccache to automatically prepend `--ccache-skip` before certain arguments.
This is at first used for fixing ccache do not recognize `icpc` AVX options. 
