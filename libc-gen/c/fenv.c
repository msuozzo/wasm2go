#include <fenv.h>

int fegetround(void) { return FE_TONEAREST; }
int fesetround(int round) { return 0; }

int feclearexcept(int mask) { return 0; }
int feraiseexcept(int mask) { return 0; }
int fetestexcept(int mask) { return 0; }

int fegetenv(fenv_t* envp) { return 0; }
int feholdexcept(fenv_t* envp) { return 0; }
int fesetenv(const fenv_t* envp) { return 0; }
int feupdateenv(const fenv_t* envp) { return 0; }

int fegetexceptflag(fexcept_t* fp, int mask) { return 0; }
int fesetexceptflag(const fexcept_t* fp, int mask) { return 0; }
