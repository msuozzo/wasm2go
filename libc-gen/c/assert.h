#undef assert

#ifdef NDEBUG
#define assert(expr) ((void)0)
#else
#define assert(expr) ((expr) ? (void)0 : __builtin_trap())
#endif

#ifndef static_assert
#define static_assert _Static_assert
#endif
