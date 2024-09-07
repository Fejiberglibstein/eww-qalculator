#define Vector(T)                                                              \
    struct {                                                                   \
        T *items;                                                              \
        uint capacity;                                                         \
        uint length;                                                           \
    }

#define vec_append(vec, item...)                                               \
    do {                                                                       \
        if ((vec)->length >= (vec)->capacity) {                                \
            (vec)->capacity =                                                  \
                ((vec)->capacity == 0) ? 1 : (vec)->capacity * 2;              \
            (vec)->items = realloc((vec)->items,                               \
                                   (vec)->capacity * sizeof(*(vec)->items));   \
        }                                                                      \
        (vec)->items[(vec)->length] = item;                                    \
        (vec)->length = (vec)->length + 1;                                     \
    } while (0)

#define vec_free(vec) free((vec)->items);
