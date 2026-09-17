#pragma once

#define FE_TONEAREST 0

int fegetround(void);
int fesetround(int);

#define FE_ALL_EXCEPT 0

typedef unsigned long fexcept_t;

int feclearexcept(int);
int feraiseexcept(int);
int fetestexcept(int);
int fegetexceptflag(fexcept_t*, int);
int fesetexceptflag(const fexcept_t*, int);

typedef struct {
  unsigned long __cw;
} fenv_t;

#define FE_DFL_ENV ((const fenv_t*)-1)

int fegetenv(fenv_t*);
int feholdexcept(fenv_t*);
int fesetenv(const fenv_t*);
int feupdateenv(const fenv_t*);
