/*
 * this is free and unencumbered software released into the public domain.
 * refer to the attached UNLICENSE or http://unlicense.org/
 * -- usage: ---------------------------------------------------------------
 * #define OPPAI_IMPLEMENTATION
 * #include "../oppai.c"
 *
 * int main() {
 *   ezpp_t ez = ezpp_new();
 *   ezpp_set_mods(ez, MODS_HD | MODS_DT);
 *   ezpp(ez, "-");
 *   printf("%gpp\n", ezpp_pp(ez));
 *   return 0;
 * }
 * ------------------------------------------------------------------------
 * $ gcc test.c
 * $ cat /path/to/file.osu | ./a.out
 */

#if defined(_WIN32) && !defined(OPPAI_IMPLEMENTATION)
#ifdef OPPAI_EXPORT
#define OPPAIAPI __declspec(dllexport)
#elif defined(OPPAI_STATIC_HEADER)
#define OPPAIAPI
#else
#define OPPAIAPI __declspec(dllimport)
#endif
#else
#define OPPAIAPI
#endif

typedef struct ezpp* ezpp_t; /* opaque handle */

OPPAIAPI ezpp_t ezpp_new(void);
OPPAIAPI void ezpp_free(ezpp_t ez);
OPPAIAPI int ezpp(ezpp_t ez, char* map);
OPPAIAPI float ezpp_pp(ezpp_t ez);
OPPAIAPI float ezpp_stars(ezpp_t ez);

/*
 * the above is all you need for basic usage. below are some advanced api's
 * and usage examples
 *
 * - if map is "-" the map is read from standard input
 * - you can use ezpp_data if you already have raw beatmap data in memory
 * - if autocalc is set to 1, the results will be automatically refreshed
 *   when you change parameters. if reparsing is required, the last passed
 *   map or map data will be used
 * - if map is 0 (NULL), difficulty calculation and map parsing are skipped
 *   and you must set at least mode, aim_stars, speed_stars, nobjects,
 *   base_ar, base_od, max_combo, nsliders, ncircles
 * - if aim_stars or speed_stars are set difficulty calculation is also
 *   skipped but values are taken from map
 * - setting mods or cs resets aim_stars and speed_stars, set those last
 * = setting end resets accuracy_percent
 * - if mode_override is set, std maps are converted to other modes
 * - mode defaults to MODE_STD or the map's mode
 * - mods default to MODS_NOMOD
 * - combo defaults to full combo
 * - nmiss defaults to 0
 * - score_version defaults to scorev1
 * - if accuracy_percent is set, n300/100/50 are automatically
 *   calculated and stored
 * - if n300/100/50 are set, accuracy_percent is automatically
 *   calculated and stored
 * - if none of the above are set, SS (100%) is assumed
 * - if end is set, the map will be cut to this object index
 * - if base_ar/od/cs are set, they will override the map's values
 * - when you change map and you're reusing the handle, you should reset
 *   ar/od/cs/hp to -1 otherwise it will override them with the previous
 *   map's values
 * - in autocalc mode, calling ezpp with a non-NULL map always resets
 *   ar/od/cs/hp overrides to -1 so you don't have to
 */

OPPAIAPI void ezpp_set_autocalc(ezpp_t ez, int autocalc);
OPPAIAPI int ezpp_autocalc(ezpp_t ez);
OPPAIAPI int ezpp_data(ezpp_t ez, char* data, int data_size);
OPPAIAPI float ezpp_aim_stars(ezpp_t ez);
OPPAIAPI float ezpp_speed_stars(ezpp_t ez);
OPPAIAPI float ezpp_aim_pp(ezpp_t ez);
OPPAIAPI float ezpp_speed_pp(ezpp_t ez);
OPPAIAPI float ezpp_acc_pp(ezpp_t ez);
OPPAIAPI float ezpp_accuracy_percent(ezpp_t ez);
OPPAIAPI int ezpp_n300(ezpp_t ez);
OPPAIAPI int ezpp_n100(ezpp_t ez);
OPPAIAPI int ezpp_n50(ezpp_t ez);
OPPAIAPI int ezpp_nmiss(ezpp_t ez);
OPPAIAPI float ezpp_ar(ezpp_t ez);
OPPAIAPI float ezpp_cs(ezpp_t ez);
OPPAIAPI float ezpp_od(ezpp_t ez);
OPPAIAPI float ezpp_hp(ezpp_t ez);
OPPAIAPI char* ezpp_artist(ezpp_t ez);
OPPAIAPI char* ezpp_artist_unicode(ezpp_t ez);
OPPAIAPI char* ezpp_title(ezpp_t ez);
OPPAIAPI char* ezpp_title_unicode(ezpp_t ez);
OPPAIAPI char* ezpp_version(ezpp_t ez);
OPPAIAPI char* ezpp_creator(ezpp_t ez);
OPPAIAPI int ezpp_ncircles(ezpp_t ez);
OPPAIAPI int ezpp_nsliders(ezpp_t ez);
OPPAIAPI int ezpp_nspinners(ezpp_t ez);
OPPAIAPI int ezpp_nobjects(ezpp_t ez);
OPPAIAPI float ezpp_odms(ezpp_t ez);
OPPAIAPI int ezpp_mode(ezpp_t ez);
OPPAIAPI int ezpp_combo(ezpp_t ez);
OPPAIAPI int ezpp_max_combo(ezpp_t ez);
OPPAIAPI int ezpp_mods(ezpp_t ez);
OPPAIAPI int ezpp_score_version(ezpp_t ez);
OPPAIAPI float ezpp_time_at(ezpp_t ez, int i);
OPPAIAPI float ezpp_strain_at(ezpp_t ez, int i, int difficulty_type);

OPPAIAPI void ezpp_set_aim_stars(ezpp_t ez, float aim_stars);
OPPAIAPI void ezpp_set_speed_stars(ezpp_t ez, float speed_stars);
OPPAIAPI void ezpp_set_base_ar(ezpp_t ez, float ar);
OPPAIAPI void ezpp_set_base_od(ezpp_t ez, float od);
OPPAIAPI void ezpp_set_base_cs(ezpp_t ez, float cs);
OPPAIAPI void ezpp_set_base_hp(ezpp_t ez, float hp);
OPPAIAPI void ezpp_set_mode_override(ezpp_t ez, int mode_override);
OPPAIAPI void ezpp_set_mode(ezpp_t ez, int mode);
OPPAIAPI void ezpp_set_mods(ezpp_t ez, int mods);
OPPAIAPI void ezpp_set_combo(ezpp_t ez, int combo);
OPPAIAPI void ezpp_set_nmiss(ezpp_t ez, int nmiss);
OPPAIAPI void ezpp_set_score_version(ezpp_t ez, int score_version);
OPPAIAPI void ezpp_set_accuracy_percent(ezpp_t ez, float accuracy_percent);
OPPAIAPI void ezpp_set_accuracy(ezpp_t ez, int n100, int n50);
OPPAIAPI void ezpp_set_end(ezpp_t ez, int end);
OPPAIAPI void ezpp_set_end_time(ezpp_t ez, float end);

/*
 * these will make a copy of mapfile/data and free it automatically. this
 * is slow but useful when working with bindings in other langs where
 * pointers to strings aren't guaranteed to persist like python3
 */
OPPAIAPI int ezpp_dup(ezpp_t ez, char* mapfile);
OPPAIAPI int ezpp_data_dup(ezpp_t ez, char* data, int data_size);

/* errors -------------------------------------------------------------- */

/*
 * all functions that return int can return errors in the form
 * of a negative value. check if the return value is < 0 and call
 * errstr to get the error message
 */

#define ERR_MORE (-1)
#define ERR_SYNTAX (-2)
#define ERR_TRUNCATED (-3)
#define ERR_NOTIMPLEMENTED (-4)
#define ERR_IO (-5)
#define ERR_FORMAT (-6)
#define ERR_OOM (-7)

OPPAIAPI char* errstr(int err);

/* version info -------------------------------------------------------- */

OPPAIAPI void oppai_version(int* major, int* minor, int* patch);
OPPAIAPI char* oppai_version_str(void);

/* --------------------------------------------------------------------- */

#define MODE_STD 0
#define MODE_TAIKO 1

#define DIFF_SPEED 0
#define DIFF_AIM 1

#define MODS_NOMOD 0
#define MODS_NF (1<<0)
#define MODS_EZ (1<<1)
#define MODS_TD (1<<2)
#define MODS_HD (1<<3)
#define MODS_HR (1<<4)
#define MODS_SD (1<<5)
#define MODS_DT (1<<6)
#define MODS_RX (1<<7)
#define MODS_HT (1<<8)
#define MODS_NC (1<<9)
#define MODS_FL (1<<10)
#define MODS_AT (1<<11)
#define MODS_SO (1<<12)
#define MODS_AP (1<<13)
#define MODS_PF (1<<14)
#define MODS_KEY4 (1<<15) /* TODO: what are these abbreviated to? */
#define MODS_KEY5 (1<<16)
#define MODS_KEY6 (1<<17)
#define MODS_KEY7 (1<<18)
#define MODS_KEY8 (1<<19)
#define MODS_FADEIN (1<<20)
#define MODS_RANDOM (1<<21)
#define MODS_CINEMA (1<<22)
#define MODS_TARGET (1<<23)
#define MODS_KEY9 (1<<24)
#define MODS_KEYCOOP (1<<25)
#define MODS_KEY1 (1<<26)
#define MODS_KEY3 (1<<27)
#define MODS_KEY2 (1<<28)
#define MODS_SCOREV2 (1<<29)
#define MODS_TOUCH_DEVICE MODS_TD
#define MODS_NOVIDEO MODS_TD /* never forget */
#define MODS_SPEED_CHANGING (MODS_DT | MODS_HT | MODS_NC)
#define MODS_MAP_CHANGING (MODS_HR | MODS_EZ | MODS_SPEED_CHANGING)

/*     this is all you need to know for normal usage. internals below    */
/* ##################################################################### */
/* ##################################################################### */
/* ##################################################################### */

#ifdef OPPAI_EXPORT
#define OPPAI_IMPLEMENTATION
#endif

#ifdef OPPAI_IMPLEMENTATION
#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>
#include <string.h>
#include <math.h>

#define OPPAI_VERSION_MAJOR 3
#define OPPAI_VERSION_MINOR 2
#define OPPAI_VERSION_PATCH 3
#define STRINGIFY_(x) #x
#define STRINGIFY(x) STRINGIFY_(x)

#define OPPAI_VERSION_STRING \
  STRINGIFY(OPPAI_VERSION_MAJOR) "." \
  STRINGIFY(OPPAI_VERSION_MINOR) "." \
  STRINGIFY(OPPAI_VERSION_PATCH)

OPPAIAPI
void oppai_version(int* major, int* minor, int* patch) {
  *major = OPPAI_VERSION_MAJOR;
  *minor = OPPAI_VERSION_MINOR;
  *patch = OPPAI_VERSION_PATCH;
}

OPPAIAPI
char* oppai_version_str() {
  return OPPAI_VERSION_STRING;
}

/* error utils --------------------------------------------------------- */

int info(char* fmt, ...) {
  int res;
  va_list va;
  va_start(va, fmt);
  res = vfprintf(stderr, fmt, va);
  va_end(va);
  return res;
}

OPPAIAPI
char* errstr(int err) {
  switch (err) {
    case ERR_MORE: return "call me again with more data";
    case ERR_SYNTAX: return "syntax error";
    case ERR_TRUNCATED:
      return "data was truncated, possibly because it was too big";
    case ERR_NOTIMPLEMENTED:
      return "requested a feature that isn't implemented";
    case ERR_IO: return "i/o error";
    case ERR_FORMAT: return "invalid input format";
    case ERR_OOM: return "out of memory";
  }
  info("W: got unknown error %d\n", err);
  return "unknown error";
}

/* math ---------------------------------------------------------------- */

#define log10f (float)log10
#define al_round(x) (float)floor((x) + 0.5f)
#define al_min(a, b) ((a) < (b) ? (a) : (b))
#define al_max(a, b) ((a) > (b) ? (a) : (b))

#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

float get_inf(void) {
  static unsigned raw = 0x7F800000;
  float* p = (float*)&raw;
  return *p;
}

float get_nan(void) {
  static unsigned raw = 0x7FFFFFFF;
  float* p = (float*)&raw;
  return *p;
}

/* dst = a - b */
void v2f_sub(float* dst, float* a, float* b) {
  dst[0] = a[0] - b[0];
  dst[1] = a[1] - b[1];
}

float v2f_len(float* v) {
  return (float)sqrt(v[0] * v[0] + v[1] * v[1]);
}

float v2f_dot(float* a, float* b) {
  return a[0] * b[0] + a[1] * b[1];
}

/* https://www.doc.ic.ac.uk/%7Eeedwards/compsys/float/nan.html */

int is_nan(float b) {
  unsigned* p = (void*)&b;
  return (
    (*p > 0x7F800000 && *p < 0x80000000) ||
    (*p > 0x7FBFFFFF && *p <= 0xFFFFFFFF)
  );
}

/* https://www.doc.ic.ac.uk/%7Eeedwards/compsys/float/nan.html */

int is_inf(float b) {
  int* p = (int*)&b;
  return *p == 0x7F800000 || *p == 0xFF800000;
}

/* string utils -------------------------------------------------------- */

int whitespace(char c) {
  switch (c) {
    case '\r':
    case '\n':
    case '\t':
    case ' ':
      return 1;
  }
  return 0;
}

/* non-null terminated string, used internally for parsing */
typedef struct slice {
  char* start;
  char* end; /* *(end - 1) is the last character */
} slice_t;

int slice_write(slice_t* s, FILE* f) {
  return (int)fwrite(s->start, 1, s->end - s->start, f);
}

int slice_whitespace(slice_t* s) {
  char* p = s->start;
  for (; p < s->end; ++p) {
    if (!whitespace(*p)) {
      return 0;
    }
  }
  return 1;
}

/* trims leading and trailing whitespace */
void slice_trim(slice_t* s) {
  for (; s->start < s->end && whitespace(*s->start); ++s->start);
  for (; s->end > s->start && whitespace(*(s->end-1)); --s->end);
}

int slice_cmp(slice_t* s, char* str) {
  int len = (int)strlen(str);
  int s_len = (int)(s->end - s->start);
  if (len < s_len) {
    return -1;
  }
  if (len > s_len) {
    return 1;
  }
  return strncmp(s->start, str, len);
}

int slice_len(slice_t* s) {
  return (int)(s->end - s->start);
}

/*
 * splits s at any of the separators in separator_list and stores
 * pointers to the strings in arr.
 * returns the number of elements written to arr.
 * if more elements than nmax are found, err is set to
 * ERR_TRUNCATED
 */
int slice_split(slice_t* s, char* separator_list, slice_t* arr,
  int nmax, int* err)
{
  int res = 0;
  char* p = s->start;
  char* pprev = p;
  if (!nmax) {
    return 0;
  }
  if (!*separator_list) {
    *arr = *s;
    return 1;
  }
  for (; p <= s->end; ++p) {
    char* sep = separator_list;
    for (; *sep; ++sep) {
      if (p >= s->end || *sep == *p) {
        if (res >= nmax) {
          *err = ERR_TRUNCATED;
          goto exit;
        }
        arr[res].start = pprev;
        arr[res].end = p;
        pprev = p + 1;
        ++res;
        break;
      }
    }
  }
exit:
  return res;
}

/* array --------------------------------------------------------------- */

#define array_t(type) \
  struct { \
    int cap; \
    int len; \
    type* data; \
  }

#define array_reserve(arr, n) \
  array_reserve_i(n, array_unpack(arr))

#define array_free(arr) \
  array_free_i(array_unpack(arr))

#define array_alloc(arr) \
  (array_reserve((arr), (arr)->len + 1) \
    ? &(arr)->data[(arr)->len++] \
    : 0)

#define array_append(arr, x) \
  (array_reserve((arr), (arr)->len + 1) \
    ? ((arr)->data[(arr)->len++] = (x), 1) \
    : 0)

/* internal helpers, not to be used directly */
#define array_unpack(arr) \
  &(arr)->cap, \
  &(arr)->len, \
  (void**)&(arr)->data, \
  (int)sizeof((arr)->data[0])

int array_reserve_i(int n, int* cap, int* len, void** data, int esize) {
  (void)len;
  if (*cap <= n) {
    void* newdata;
    int newcap = *cap ? *cap * 2 : 16;
    newdata = realloc(*data, esize * newcap);
    if (!newdata) {
      return 0;
    }
    *data = newdata;
    *cap = newcap;
  }
  return 1;
}

void array_free_i(int* cap, int* len, void** data, int esize) {
  (void)esize;
  free(*data);
  *cap = 0;
  *len = 0;
  *data = 0;
}

/* --------------------------------------------------------------------- */

#define OBJ_CIRCLE (1<<0)
#define OBJ_SLIDER (1<<1)
#define OBJ_SPINNER (1<<3)

#define SOUND_NONE 0
#define SOUND_NORMAL (1<<0)
#define SOUND_WHISTLE (1<<1)
#define SOUND_FINISH (1<<2)
#define SOUND_CLAP (1<<3)

typedef struct timing {
  float time;        /* milliseconds */
  float ms_per_beat;
  int change;        /* if 0, ms_per_beat is -100.0f * sv_multiplier */
  float px_per_beat;

  /* taiko stuff */
  float beat_len;
  float velocity;
} timing_t;

typedef struct object {
  float time; /* milliseconds */
  int type;

  /* only for taiko maps */
  int nsound_types;
  int* sound_types;

  /* only used by d_calc */
  float normpos[2];
  float angle;
  float strains[2];
  int is_single; /* 1 if diff calc sees this as a singletap */
  float delta_time;
  float d_distance;
  int timing_point;

  float pos[2];
  float distance;  /* only for sliders */
  int repetitions;

  /* taiko stuff */
  float duration;
  float tick_spacing;
  int slider_is_drum_roll;
} object_t;

/*
 * exposing the struct would cut down lines of code but makes it harder
 * to use from langs that aren't c/c++ or don't have the same memory
 * alignment etc
 */

#define AUTOCALC_BIT (1<<0)
#define OWNS_MAP_BIT (1<<1) /* map/data freed on ezpp{,_data}, ezpp_free */

struct ezpp {
  char* map;
  char* data;
  int data_size;
  int flags;
  int format_version;
  int mode, mode_override, original_mode;
  int score_version;
  int mods, combo;
  float accuracy_percent;
  int n300, n100, n50, nmiss;
  int end;
  float end_time;
  float base_ar, base_cs, base_od, base_hp;
  int max_combo;
  char* title;
  char* title_unicode;
  char* artist;
  char* artist_unicode;
  char* creator;
  char* version;
  int ncircles, nsliders, nspinners, nobjects;
  float ar, od, cs, hp, odms, sv, tick_rate, speed_mul;
  float stars;
  float aim_stars, aim_difficulty, aim_length_bonus;
  float speed_stars, speed_difficulty, speed_length_bonus;
  float pp, aim_pp, speed_pp, acc_pp;

  /* parser */
  char section[64];
  char buf[0xFFFF];
  int p_flags;
  array_t(object_t) objects;
  array_t(timing_t) timing_points;

  /* diffcalc */
  float interval_end;
  float max_strain;
  array_t(float) highest_strains;

  /* allocator */
  char* block;
  char* end_of_block;
  array_t(char*) blocks;
};

/* memory arena (allocator) -------------------------------------------- */

#define M_ALIGN sizeof(void*)
#define M_BLOCK_SIZE 4096

/* aligns x down to a power-of-two value a */
#define bit_align_down(x, a) \
  ((x) & ~((a) - 1))

/* aligns x up to a power-of-two value a */
#define bit_align_up(x, a) \
  bit_align_down((x) + (a) - 1, a)

int m_reserve(ezpp_t ez, int min_size) {
  int size;
  char* new_block;
  if (ez->end_of_block - ez->block >= min_size) {
    return 1;
  }
  size = bit_align_up(al_max(min_size, M_BLOCK_SIZE), M_ALIGN);
  new_block = malloc(size);
  if (!new_block) {
    return 0;
  }
  ez->block = new_block;
  ez->end_of_block = new_block + size;
  array_append(&ez->blocks, ez->block);
  return 1;
}

void* m_alloc(ezpp_t ez, int size) {
  void* res;
  if (!m_reserve(ez, size)) {
    return 0;
  }
  size = bit_align_up(size, M_ALIGN);
  res = ez->block;
  ez->block += size;
  return res;
}

char* m_strndup(ezpp_t ez, char* s, int n) {
  char* res = m_alloc(ez, n + 1);
  if (res) {
    memcpy(res, s, n);
    res[n] = 0;
  }
  return res;
}

void m_free(ezpp_t ez) {
  int i;
  for (i = 0; i < ez->blocks.len; ++i) {
    free(ez->blocks.data[i]);
  }
  array_free(&ez->blocks);
  ez->block = 0;
  ez->end_of_block = 0;
}

/* mods ---------------------------------------------------------------- */

float od10_ms[] = { 20, 20 }; /* std, taiko */
float od0_ms[] = { 80, 50 };
#define AR0_MS 1800.0f
#define AR5_MS 1200.0f
#define AR10_MS 450.0f

float od_ms_step[] = { 6.0f, 3.0f };
#define AR_MS_STEP1 120.f /* ar0-5 */
#define AR_MS_STEP2 150.f /* ar5-10 */

/*
 * stats must be capped to 0-10 before HT/DT which brings them to a range
 * of -4.42f to 11.08f for OD and -5 to 11 for AR
 */

int mods_apply(ezpp_t ez) {
  float od_ar_hp_multiplier, cs_multiplier, arms;

  switch (ez->mode) {
  case MODE_STD:
  case MODE_TAIKO:
    break;
  default:
    info("this gamemode is not yet supported for mods calc\n");
    return ERR_NOTIMPLEMENTED;
  }

  ez->speed_mul = 1;

  if (!(ez->mods & MODS_MAP_CHANGING)) {
    ez->odms = od0_ms[ez->mode] - (float)ceil(od_ms_step[ez->mode] * ez->od);
    return 0;
  }

  if (ez->mods & (MODS_DT | MODS_NC)) {
    ez->speed_mul *= 1.5f;
  }
  if (ez->mods & MODS_HT) {
    ez->speed_mul *= 0.75f;
  }

  /* global multipliers */
  od_ar_hp_multiplier = 1;
  if (ez->mods & MODS_HR) od_ar_hp_multiplier *= 1.4f;
  if (ez->mods & MODS_EZ) od_ar_hp_multiplier *= 0.5f;

  ez->od *= od_ar_hp_multiplier;
  ez->odms = od0_ms[ez->mode] - (float)ceil(od_ms_step[ez->mode] * ez->od);
  ez->odms = al_min(od0_ms[ez->mode], al_max(od10_ms[ez->mode], ez->odms));
  ez->odms /= ez->speed_mul;
  ez->od = (od0_ms[ez->mode] - ez->odms) / od_ms_step[ez->mode];

  ez->ar *= od_ar_hp_multiplier;
  arms = ez->ar <= 5
    ? (AR0_MS - AR_MS_STEP1 * (ez->ar - 0))
    : (AR5_MS - AR_MS_STEP2 * (ez->ar - 5));
  arms = al_min(AR0_MS, al_max(AR10_MS, arms));
  arms /= ez->speed_mul;
  ez->ar = arms > AR5_MS
    ? (0 + (AR0_MS - arms) / AR_MS_STEP1)
    : (5 + (AR5_MS - arms) / AR_MS_STEP2);

  cs_multiplier = 1;
  if (ez->mods & MODS_HR) cs_multiplier = 1.3f;
  if (ez->mods & MODS_EZ) cs_multiplier = 0.5f;
  ez->cs *= cs_multiplier;
  ez->cs = al_max(0.0f, al_min(10.0f, ez->cs));

  ez->hp = al_min(ez->hp * od_ar_hp_multiplier, 10);

  return 0;
}

/* beatmap parser ------------------------------------------------------ */

/*
 * comments in beatmaps can only be an entire line because
 * some properties such as author can contain //
 *
 * all p_* functions expect s to be a single line and trimmed
 * on errors, p_* functions return < 0 error codes otherwise they
 * return n bytes consumed
 */

#define P_OVERRIDE_MODE (1<<0) /* mode_override */
#define P_FOUND_AR (1<<1)

#define CIRCLESIZE_BUFF_TRESHOLD 30.0f /* non-normalized diameter */
#define PLAYFIELD_WIDTH 512.0f /* in osu!pixels */
#define PLAYFIELD_HEIGHT 384.0f

float playfield_center[] = {
  PLAYFIELD_WIDTH / 2.0f, PLAYFIELD_HEIGHT / 2.0f
};

void print_line(slice_t* line) {
  info("in line: ");
  slice_write(line, stderr);
  info("\n");
}

int p_warn(char* e, slice_t* line) {
  info(e);
  info("\n");
  print_line(line);
  return 0;
}

/* consume until any of the characters in separators is found */
int p_consume_til(slice_t* s, char* separators, slice_t* dst) {
  char* p = s->start;
  dst->start = s->start;
  for (; p < s->end; ++p) {
    char* sep;
    for (sep = separators; *sep; ++sep) {
      if (*p == *sep) {
        dst->start = s->start;
        dst->end = p;
        return (int)(p - s->start);
      }
    }
  }
  dst->end = p;
  return ERR_MORE;
}

float p_float(slice_t* value) {
  float res;
  char* p = value->start;
  if (*p == '-') {
    res = -1;
    ++p;
  } else {
    res = 1;
  }
  /* infinity symbol */
  if (!strncmp(p, "\xe2\x88\x9e", 3)) {
    res *= get_inf();
  } else {
    if (sscanf(value->start, "%f", &res) != 1) {
      info("W: failed to parse float ");
      slice_write(value, stderr);
      info("\n");
      res = 0;
    }
  }
  return res;
}

/* [name] */
int p_section_name(slice_t* s, slice_t* name) {
  int n;
  slice_t p = *s;
  if (*p.start++ != '[') {
    return ERR_SYNTAX;
  }
  n = p_consume_til(&p, "]", name);
  if (n < 0) {
    return n;
  }
  p.start += n;
  if (p.start != p.end - 1) { /* must end in ] */
    return ERR_SYNTAX;
  }
  return (int)(p.start - s->start);
}

/* name: value (results are trimmed) */
int p_property(slice_t* s, slice_t* name, slice_t* value) {
  int n;
  char* p = s->start;
  n = p_consume_til(s, ":", name);
  if (n < 0) {
    return n;
  }
  p += n;
  ++p; /* skip : */
  value->start = p;
  value->end = s->end;
  slice_trim(name);
  slice_trim(value);
  return (int)(s->end - s->start);
}

char* p_slicedup(ezpp_t ez, slice_t* s) {
  return m_strndup(ez, s->start, slice_len(s));
}

int p_metadata(ezpp_t ez, slice_t* line) {
  slice_t name, value;
  int n = p_property(line, &name, &value);
  if (n < 0) {
    return p_warn("W: malformed metadata line", line);
  }
  if (!slice_cmp(&name, "Title")) {
    ez->title = p_slicedup(ez, &value);
  } else if (!slice_cmp(&name, "TitleUnicode")) {
    ez->title_unicode = p_slicedup(ez, &value);
  } else if (!slice_cmp(&name, "Artist")) {
    ez->artist = p_slicedup(ez, &value);
  } else if (!slice_cmp(&name, "ArtistUnicode")) {
    ez->artist_unicode = p_slicedup(ez, &value);
  } else if (!slice_cmp(&name, "Creator")) {
    ez->creator = p_slicedup(ez, &value);
  } else if (!slice_cmp(&name, "Version")) {
    ez->version = p_slicedup(ez, &value);
  }
  return n;
}

int p_general(ezpp_t ez, slice_t* line) {
  slice_t name, value;
  int n;
  n = p_property(line, &name, &value);
  if (n < 0) {
    return p_warn("W: malformed general line", line);
  }
  if (!slice_cmp(&name, "Mode")) {
    if (sscanf(value.start, "%d", &ez->original_mode) != 1){
      return ERR_SYNTAX;
    }
    if (ez->p_flags & P_OVERRIDE_MODE) {
      ez->mode = ez->mode_override;
    } else {
      ez->mode = ez->original_mode;
    }
    switch (ez->mode) {
    case MODE_STD:
    case MODE_TAIKO:
      break;
    default:
      return ERR_NOTIMPLEMENTED;
    }
  }
  return n;
}

int p_difficulty(ezpp_t ez, slice_t* line) {
  slice_t name, value;
  int n = p_property(line, &name, &value);
  if (n < 0) {
    return p_warn("W: malformed difficulty line", line);
  }
  if (!slice_cmp(&name, "CircleSize")) {
    ez->cs = p_float(&value);
  } else if (!slice_cmp(&name, "OverallDifficulty")) {
    ez->od = p_float(&value);
  } else if (!slice_cmp(&name, "ApproachRate")) {
    ez->ar = p_float(&value);
    ez->p_flags |= P_FOUND_AR;
  } else if (!slice_cmp(&name, "HPDrainRate")) {
    ez->hp = p_float(&value);
  } else if (!slice_cmp(&name, "SliderMultiplier")) {
    ez->sv = p_float(&value);
  } else if (!slice_cmp(&name, "SliderTickRate")) {
    ez->tick_rate = p_float(&value);
  }
  return n;
}

/*
 * time, ms_per_beat, time_signature_id, sample_set_id,
 * sample_bank_id, sample_volume, is_timing_change, effect_flags
 *
 * everything after ms_per_beat is optional
 */
int p_timing(ezpp_t ez, slice_t* line) {
  int res = 0;
  int n, i;
  int err = 0;
  slice_t split[8];
  timing_t* t = array_alloc(&ez->timing_points);

  if (!t) {
    return ERR_OOM;
  }

  t->change = 1;

  n = slice_split(line, ",", split, 8, &err);
  if (err < 0) {
    if (err == ERR_TRUNCATED) {
      info("W: timing point with trailing values");
      print_line(line);
    } else {
      return err;
    }
  }

  if (n < 2) {
    return p_warn("W: malformed timing point", line);
  }

  res = (int)(split[n - 1].end - line->start);
  for (i = 0; i < n; ++i) {
    slice_trim(&split[i]);
  }

  t->time = p_float(&split[0]);
  t->ms_per_beat = p_float(&split[1]);

  if (n >= 7) {
    if (slice_len(&split[6]) < 1) {
      t->change = 1;
    } else {
      t->change = *split[6].start != '0';
    }
  }

  return res;
}

int p_objects(ezpp_t ez, slice_t* line) {
  object_t* o;
  int err = 0;
  int ne;
  slice_t e[11];

  if (ez->end > 0 && ez->objects.len >= ez->end) {
    return 0;
  }

  o = array_alloc(&ez->objects);

  if (o) {
    memset(o, 0, sizeof(*o));
  } else {
    return ERR_OOM;
  }

  ne = slice_split(line, ",", e, 11, &err);
  if (err < 0) {
    if (err == ERR_TRUNCATED) {
      info("W: object with trailing values\n");
      print_line(line);
    } else {
      return err;
    }
  }

  if (ne < 5) {
    return p_warn("W: malformed hitobject", line);
  }

  o->time = p_float(&e[2]);
  if (is_inf(o->time)) {
    o->time = 0.0f;
    info("W: object with infinite time\n");
    print_line(line);
  }

  if (ez->end_time > 0 && o->time >= ez->end_time) {
    --ez->objects.len;
    return 0;
  }

  if (sscanf(e[3].start, "%d", &o->type) != 1) {
    p_warn("W: malformed hitobject type", line);
    o->type = OBJ_CIRCLE;
  }

  if (ez->mode == MODE_TAIKO) {
    int* sound_type = m_alloc(ez, sizeof(int));
    if (!sound_type) {
      return ERR_OOM;
    }
    if (sscanf(e[4].start, "%d", sound_type) != 1) {
      p_warn("W: malformed hitobject sound type", line);
      *sound_type = SOUND_NORMAL;
    }
    o->nsound_types = 1;
    o->sound_types = sound_type;
    /* wastes 4 bytes when you have per-node sounds but w/e */
  }

  if (o->type & OBJ_CIRCLE) {
    ++ez->ncircles;
    o->pos[0] = p_float(&e[0]);
    o->pos[1] = p_float(&e[1]);
  }

  /* ?,?,?,?,?,end_time,custom_sample_banks */
  else if (o->type & OBJ_SPINNER) {
    ++ez->nspinners;
  }

  /*
   * x,y,time,type,sound_type,points,repetitions,distance,
   * per_node_sounds,per_node_samples,custom_sample_banks
   */
  else if (o->type & OBJ_SLIDER) {
    ++ez->nsliders;
    if (ne < 7) {
      return p_warn("W: malformed slider", line);
    }

    o->pos[0] = p_float(&e[0]);
    o->pos[1] = p_float(&e[1]);

    if (sscanf(e[6].start, "%d", &o->repetitions) != 1) {
      o->repetitions = 1;
      p_warn("W: malformed slider repetitions", line);
    }

    if (ne > 7) {
      o->distance = p_float(&e[7]);
    } else {
      o->distance = 0;
    }

    /* per-node sound types */
    if (ez->mode == MODE_TAIKO && ne > 8 && slice_len(&e[8]) > 0) {
      slice_t p = e[8];
      int i, nodes;
      int sound_type = o->sound_types[0];

      /*
       * TODO: there's probably something subtly wrong with this.
       * sometimes we get less sound types than nodes
       * also I don't know if I'm supposed to include the previous
       * sound type from the single sound_type field
       */

      /* repeats + head and tail. no repeats is 0 repetition */
      nodes = o->repetitions + 1;
      if (nodes < 0 || nodes > 1000) {
        /* TODO: not sure if 1000 limit is enough */
        p_warn("W: malformed node count", line);
        return ERR_SYNTAX;
      }
      o->sound_types = m_alloc(ez, sizeof(int) * nodes);
      if (!o->sound_types) {
        return ERR_OOM;
      }

      o->nsound_types = nodes;
      for (i = 0; i < nodes; ++i) {
        o->sound_types[i] = sound_type;
      }

      for (i = 0; i < nodes; ++i) {
        slice_t node;
        int n;
        int type;
        node.start = node.end = 0;
        n = p_consume_til(&p, "|", &node);
        if (n < 0 && n != ERR_MORE) {
          return n;
        }
        if (node.start >= node.end || !node.start) {
          break;
        }
        p.start += n + 1;
        if (sscanf(node.start, "%d", &type) != 1) {
          p_warn("W: malformed sound type", line);
          break;
        }
        o->sound_types[i] = type;
        if (p.start >= p.end) {
          break;
        }
      }
    }
  }

  return (int)(e[ne - 1].end - line->start);
}

int p_line(ezpp_t ez, slice_t* line) {
  int n = 0;

  if (line->start >= line->end) {
    /* empty line */
    return 0;
  }

  if (slice_whitespace(line)) {
    return (int)(line->end - line->start);
  }

  /* comments (according to lazer) */
  switch (*line->start) {
    case ' ':
    case '_':
      return (int)(line->end - line->start);
  }

  /* from here on we don't care about leading or trailing whitespace */
  slice_trim(line);

  /* C++ style comments */
  if (!strncmp(line->start, "//", 2)) {
    return 0;
  }

  /* new section */
  if (*line->start == '[') {
    slice_t section;
    int len;
    n = p_section_name(line, &section);
    if (n < 0) {
      return n;
    }
    if ((int)(section.end - section.start) >= (int)sizeof(ez->section)) {
      p_warn("W: truncated long section name", line);
    }
    len = al_min((int)sizeof(ez->section) - 1,
      (int)(section.end - section.start));
    memcpy(ez->section, section.start, len);
    ez->section[len] = 0;
    return n;
  }

  if (!strcmp(ez->section, "Metadata")) {
    n = p_metadata(ez, line);
  } else if (!strcmp(ez->section, "General")) {
    n = p_general(ez, line);
  } else if (!strcmp(ez->section, "Difficulty")) {
    n = p_difficulty(ez, line);
  } else if (!strcmp(ez->section, "TimingPoints")) {
    n = p_timing(ez, line);
  } else if (!strcmp(ez->section, "HitObjects")) {
    n = p_objects(ez, line);
  } else {
    char* p = line->start;
    char* fmt_str = "file format v";
    for (; p < line->end && strncmp(p, fmt_str, 13); ++p);
    p += 13;
    if (p < line->end) {
      if (sscanf(p, "%d", &ez->format_version) == 1) {
        return (int)(line->end - line->start);
      }
    }
  }

  return n;
}

void p_end(ezpp_t ez) {
  int i;
  float infinity = get_inf();
  float tnext = -infinity;
  int tindex = -1;
  float ms_per_beat = infinity;
  float radius, scaling_factor;
  float legacy_multiplier = 1;

  if (!(ez->p_flags & P_FOUND_AR)) {
    /* in old maps ar = od */
    ez->ar = ez->od;
  }

  if (!ez->title_unicode) {
    ez->title_unicode = ez->title;
  }

  if (!ez->artist_unicode) {
    ez->artist_unicode = ez->artist;
  }

  #define s(x) ez->x = ez->x ? ez->x : "(null)"
  s(title);
  s(title_unicode);
  s(artist);
  s(artist_unicode);
  s(creator);
  s(version);
  #undef s

  if (ez->base_ar < 0) ez->base_ar = ez->ar;
  else ez->ar = ez->base_ar;
  if (ez->base_cs < 0) ez->base_cs = ez->cs;
  else ez->cs = ez->base_cs;
  if (ez->base_od < 0) ez->base_od = ez->od;
  else ez->od = ez->base_od;
  if (ez->base_hp < 0) ez->base_hp = ez->hp;
  else ez->hp = ez->base_hp;
  mods_apply(ez);

  if (ez->mode == MODE_TAIKO && ez->mode != ez->original_mode) {
    legacy_multiplier = 1.4f;
    ez->sv *= legacy_multiplier;
  }

  for (i = 0; i < ez->timing_points.len; ++i) {
    timing_t* t = &ez->timing_points.data[i];
    float sv_multiplier = 1.0f;
    if (t->change) {
      ms_per_beat = t->ms_per_beat;
    }
    if (!t->change && t->ms_per_beat < 0) {
      sv_multiplier = -100.0f / t->ms_per_beat;
    }
    t->beat_len = ms_per_beat / sv_multiplier;
    t->px_per_beat = ez->sv * 100.0f;
    t->velocity = 100.0f * ez->sv / t->beat_len;
    if (ez->format_version >= 8) {
      t->beat_len *= sv_multiplier;
      t->px_per_beat *= sv_multiplier;
    }
  }

  /*
   * sliders get 2 + ticks combo (head, tail and ticks) each repetition adds
   * an extra combo and an extra set of ticks
   *
   * calculate the number of slider ticks for one repetition
   * ---
   * example: a 3.75f beats slider at 1x tick rate will go:
   * beat0 (head), beat1 (tick), beat2(tick), beat3(tick),
   * beat3.75f(tail)
   * so all we have to do is ceil the number of beats and subtract 1 to take
   * out the tail
   * ---
   * the -0.1f is there to prevent ceil from ceiling whole values like 1.0f to
   * 2.0f randomly
   */

  ez->nobjects = ez->objects.len;
  ez->max_combo = ez->nobjects;

  /* spinners don't give combo in taiko */
  if (ez->mode == MODE_TAIKO) {
    ez->max_combo -= ez->nspinners + ez->nsliders;
  }

  /*
   * positions are normalized on circle radius so that we
   * can calc as if everything was the same circlesize
   * this should really be in diffcalc functions but putting it here
   * makes it so that i only traverse the hitobjects twice in total
   */
  radius = (
    (PLAYFIELD_WIDTH / 16.0f) *
    (1.0f - 0.7f * ((float)ez->cs - 5.0f) / 5.0f)
  );

  scaling_factor = 52.0f / radius;

  /* cs buff (originally from osuElements) */
  if (radius < CIRCLESIZE_BUFF_TRESHOLD) {
    scaling_factor *=
      1.0f + al_min((CIRCLESIZE_BUFF_TRESHOLD - radius), 5.0f) / 50.0f;
  }

  for (i = 0; i < ez->objects.len; ++i) {
    object_t* o = &ez->objects.data[i];
    timing_t* t;
    int ticks;
    float num_beats;
    float* pos;
    float dot, det;

    if (o->type & OBJ_SPINNER) {
      pos = playfield_center;
    } else {
      /* sliders also begin with pos so it's fine */
      pos = o->pos;
    }
    o->normpos[0] = pos[0] * scaling_factor;
    o->normpos[1] = pos[1] * scaling_factor;

    /* angle data */
    if (i >= 2) {
      object_t* prev1 = &ez->objects.data[i - 1];
      object_t* prev2 = &ez->objects.data[i - 2];
      float v1[2], v2[2];
      v2f_sub(v1, prev2->normpos, prev1->normpos);
      v2f_sub(v2, o->normpos, prev1->normpos);
      dot = v2f_dot(v1, v2);
      det = v1[0] * v2[1] - v1[1] * v2[0];
      o->angle = (float)fabs(atan2(det, dot));
    } else {
      o->angle = get_nan();
    }

    /* keep track of the current timing point */
    while (o->time >= tnext) {
      ++tindex;
      if (tindex + 1 < ez->timing_points.len) {
        tnext = ez->timing_points.data[tindex + 1].time;
      } else {
        tnext = infinity;
      }
    }
    o->timing_point = tindex;
    t = &ez->timing_points.data[tindex];

    if (!(o->type & OBJ_SLIDER)) {
      continue;
    }

    o->duration = o->distance * o->repetitions / t->velocity;
    o->duration *= legacy_multiplier;
    o->tick_spacing = al_min(t->beat_len / ez->tick_rate,
        o->duration / o->repetitions);
    o->slider_is_drum_roll = (
      o->tick_spacing > 0 && o->duration < 2 * t->beat_len
    );

    /* slider ticks for max_combo */
    switch (ez->mode) {
      case MODE_TAIKO: {
        if (o->slider_is_drum_roll && ez->mode != ez->original_mode) {
          ez->max_combo += (int)(
            ceil((o->duration + o->tick_spacing / 8) / o->tick_spacing)
          );
        }
        break;
      }
      case MODE_STD:
        /* std slider ticks */
        num_beats = (o->distance * o->repetitions) / t->px_per_beat;

        ticks = (int)ceil(
          (num_beats - 0.1f) / o->repetitions * ez->tick_rate
        );
        --ticks;

        ticks *= o->repetitions;     /* account for repetitions */
        ticks += o->repetitions + 1; /* add heads and tails */

        /*
         * actually doesn't include first head because we already
         * added it by setting res = nobjects
         */
        ez->max_combo += al_max(0, ticks - 1);
        break;
    }
  }
}

void p_reset(ezpp_t ez) {
  ez->ncircles = ez->nsliders = ez->nspinners = ez->nobjects = 0;
  ez->objects.len = 0;
  ez->timing_points.len = 0;
  m_free(ez);
  memset(ez->section, 0, sizeof(ez->section));
}

int p_map(ezpp_t ez, FILE* f) {
  char* p;
  char* bufend;
  int c, n;
  slice_t line;
  if (!f) {
    return ERR_IO;
  }

  p_reset(ez);

  /* reading loop */
  bufend = ez->buf + sizeof(ez->buf) - 1;
  do {
    p = ez->buf;
    for (;;) {
      if (p >= bufend) {
        return ERR_TRUNCATED;
      }
      c = fgetc(f);
      if (c == '\n' || c == EOF) {
        break;
      }
      *p++ = (char)c;
    }
    *p = 0;
    line.start = ez->buf;
    line.end = p;
    n = p_line(ez, &line);
    if (n < 0) {
      return n;
    }
  } while (c != EOF);

  p_end(ez);
  ez->nobjects = ez->objects.len;

  return (int)(p - ez->buf);
}

/* TODO: see if i can shrink this function */
int p_map_mem(ezpp_t ez, char* data, int data_size) {
  int res = 0;
  int n;
  int nlines = 0; /* complete lines in the current chunk */
  slice_t s; /* points to the remaining data in buf */

  if (!data || data_size == 0) {
    return ERR_IO;
  }

  p_reset(ez);

  s.start = data;
  s.end = data + data_size;

  /* parsing loop */
  for (; s.start < s.end; ) {
    slice_t line;
    n = p_consume_til(&s, "\n", &line);

    if (n < 0) {
      if (n != ERR_MORE) {
        return n;
      }
      if (!nlines) {
        /* line doesn't fit the entire buffer */
        return ERR_TRUNCATED;
      }
      /* EOF, so we must process the remaining data as a line */
      line = s;
      n = (int)(s.end - s.start);
    } else {
      ++n; /* also skip the \n */
    }

    res += n;
    s.start += n;
    ++nlines;

    n = p_line(ez, &line);
    if (n < 0) {
      return n;
    }

    res += n;
  }

  p_end(ez);

  return res;
}

/* diff calc ----------------------------------------------------------- */

/* based on tom94's osu!tp aimod and osuElements */

#define SINGLE_SPACING 125.0f
#define STAR_SCALING_FACTOR 0.0675f /* star rating multiplier */
#define EXTREME_SCALING_FACTOR 0.5f /* used to mix aim/speed stars */
#define STRAIN_STEP 400.0f /* diffcalc uses peak strains of 400ms chunks */
#define DECAY_WEIGHT 0.9f /* peak strains are added in a weighed sum */
#define MAX_SPEED_BONUS 45.0f /* ~330BPM 1/4 streams */
#define MIN_SPEED_BONUS 75.0f /* ~200BPM 1/4 streams */
#define ANGLE_BONUS_SCALE 90
#define AIM_TIMING_THRESHOLD 107
#define SPEED_ANGLE_BONUS_BEGIN (5 * M_PI / 6)
#define AIM_ANGLE_BONUS_BEGIN (M_PI / 3)
float decay_base[] = { 0.3f, 0.15f }; /* strains decay per interval */
float weight_scaling[] = { 1400.0f, 26.25f }; /* balances aim/speed */

/*
 * TODO: unbloat these params
 * this function has become a mess with the latest changes, I should split
 * it into separate funcs for speed and im
 */
float d_spacing_weight(float distance, float delta_time,
  float prev_distance, float prev_delta_time,
  float angle, int type, int* is_single)
{
  float angle_bonus;
  float strain_time = al_max(delta_time, 50.0f);
  switch (type) {
    case DIFF_SPEED: {
      float speed_bonus;
      *is_single = distance > SINGLE_SPACING;
      distance = al_min(distance, SINGLE_SPACING);
      delta_time = al_max(delta_time, MAX_SPEED_BONUS);
      speed_bonus = 1.0f;
      if (delta_time < MIN_SPEED_BONUS) {
        speed_bonus += (float)
          pow((MIN_SPEED_BONUS - delta_time) / 40.0f, 2);
      }
      angle_bonus = 1.0f;
      if (!is_nan(angle) && angle < SPEED_ANGLE_BONUS_BEGIN) {
        float s = (float)sin(1.5 * (SPEED_ANGLE_BONUS_BEGIN - angle));
        angle_bonus += (float)pow(s, 2) / 3.57f;
        if (angle < M_PI / 2) {
          angle_bonus = 1.28f;
          if (distance < ANGLE_BONUS_SCALE && angle < M_PI / 4) {
            angle_bonus += (1 - angle_bonus)
              * al_min((ANGLE_BONUS_SCALE - distance) / 10, 1);
          }
          else if (distance < ANGLE_BONUS_SCALE) {
            angle_bonus += (1 - angle_bonus)
              * al_min((ANGLE_BONUS_SCALE - distance) / 10, 1)
              * (float)sin((M_PI / 2 - angle) * 4 / M_PI);
          }
        }
      }
      return (
        (1 + (speed_bonus - 1) * 0.75f) *
        angle_bonus *
        (0.95f + speed_bonus * (float)pow(distance / SINGLE_SPACING, 3.5))
      ) / strain_time;
    }
    case DIFF_AIM: {
      float result = 0;
      float weighted_distance;
      float prev_strain_time = al_max(prev_delta_time, 50.0f);
      if (!is_nan(angle) && angle > AIM_ANGLE_BONUS_BEGIN) {
        angle_bonus = (float)sqrt(
          al_max(prev_distance - ANGLE_BONUS_SCALE, 0)
          * pow(sin(angle - AIM_ANGLE_BONUS_BEGIN), 2)
          * al_max(distance - ANGLE_BONUS_SCALE, 0)
        );
        result = 1.5f * (float)pow(al_max(0, angle_bonus), 0.99)
          / al_max(AIM_TIMING_THRESHOLD, prev_strain_time);
      }
      weighted_distance = (float)pow(distance, 0.99);
      return al_max(
        result + weighted_distance /
          al_max(AIM_TIMING_THRESHOLD, strain_time),
        weighted_distance / strain_time
      );
    }
  }
  return 0.0f;
}

void d_calc_strain(int type, object_t* o, object_t* prev, float speedmul) {
  float res = 0;
  float time_elapsed = (o->time - prev->time) / speedmul;
  float decay = (float)pow(decay_base[type], time_elapsed / 1000.0f);
  float scaling = weight_scaling[type];

  o->delta_time = time_elapsed;

  /* this implementation doesn't account for sliders */
  if (o->type & (OBJ_SLIDER | OBJ_CIRCLE)) {
    float diff[2];
    v2f_sub(diff, o->normpos, prev->normpos);
    o->d_distance = v2f_len(diff);
    res = d_spacing_weight(o->d_distance, time_elapsed, prev->d_distance,
      prev->delta_time, o->angle, type, &o->is_single);
    res *= scaling;
  }

  o->strains[type] = prev->strains[type] * decay + res;
}

#if defined(_WIN32)
#define FORCECDECL __cdecl
#else
#define FORCECDECL
#endif

int FORCECDECL dbl_desc(void const* a, void const* b) {
  float x = *(float const*)a;
  float y = *(float const*)b;
  if (x < y) return 1;
  if (x == y) return 0;
  return -1;
}

int d_update_max_strains(ezpp_t ez, float decay_factor,
  float cur_time, float prev_time, float cur_strain, float prev_strain,
  int first_obj)
{
  /* make previous peak strain decay until the current obj */
  while (cur_time > ez->interval_end) {
    if (!array_append(&ez->highest_strains, ez->max_strain)) {
      return ERR_OOM;
    }
    if (first_obj) {
      ez->max_strain = 0;
    } else {
      float decay;
      decay = (float)pow(decay_factor,
        (ez->interval_end - prev_time) / 1000.0f);
      ez->max_strain = prev_strain * decay;
    }
    ez->interval_end += STRAIN_STEP * ez->speed_mul;
  }

  ez->max_strain = al_max(ez->max_strain, cur_strain);
  return 0;
}

void d_weigh_strains(ezpp_t ez, float* pdiff, float* ptotal) {
  int i;
  int nstrains = 0;
  float* strains;
  float total = 0;
  float difficulty = 0;
  float weight = 1.0f;

  strains = (float*)ez->highest_strains.data;
  nstrains = ez->highest_strains.len;

  /* sort strains from highest to lowest */
  qsort(strains, nstrains, sizeof(float), dbl_desc);

  for (i = 0; i < nstrains; ++i) {
    total += (float)pow(strains[i], 1.2);
    difficulty += strains[i] * weight;
    weight *= DECAY_WEIGHT;
  }

  *pdiff = difficulty;
  if (ptotal) {
    *ptotal = total;
  }
}

int d_calc_individual(ezpp_t ez, int type) {
  int i;

  /*
   * the first object doesn't generate a strain,
   * so we begin with an incremented interval end
   */
  ez->max_strain = 0.0f;
  ez->interval_end = (
    (float)ceil(ez->objects.data[0].time / (STRAIN_STEP * ez->speed_mul))
    * STRAIN_STEP * ez->speed_mul
  );
  ez->highest_strains.len = 0;

  for (i = 0; i < ez->objects.len; ++i) {
    int err;
    object_t* o = &ez->objects.data[i];
    object_t* prev = 0;
    float prev_time = 0, prev_strain = 0;
    if (i > 0) {
      prev = &ez->objects.data[i - 1];
      d_calc_strain(type, o, prev, ez->speed_mul);
      prev_time = prev->time;
      prev_strain = prev->strains[type];
    }
    err = d_update_max_strains(ez, decay_base[type], o->time, prev_time,
      o->strains[type], prev_strain, i == 0);
    if (err < 0) {
      return err;
    }
  }

  /*
   * the peak strain will not be saved for
   * the last section in the above loop
   */
  if (!array_append(&ez->highest_strains, ez->max_strain)) {
    return ERR_OOM;
  }

  switch (type) {
    case DIFF_SPEED:
      d_weigh_strains(ez, &ez->speed_stars, &ez->speed_difficulty);
      break;
    case DIFF_AIM:
      d_weigh_strains(ez, &ez->aim_stars, &ez->aim_difficulty);
      break;
  }
  return 0;
}

float d_length_bonus(float stars, float difficulty) {
  return 0.32f + 0.5f * (log10f(difficulty + stars) - log10f(stars));
}

int d_std(ezpp_t ez) {
  int res;

  res = d_calc_individual(ez, DIFF_SPEED);
  if (res < 0) {
    return res;
  }

  res = d_calc_individual(ez, DIFF_AIM);
  if (res < 0) {
    return res;
  }

  ez->aim_length_bonus = d_length_bonus(ez->aim_stars, ez->aim_difficulty);
  ez->speed_length_bonus = d_length_bonus(ez->speed_stars, ez->speed_difficulty);
  ez->aim_stars = (float)sqrt(ez->aim_stars) * STAR_SCALING_FACTOR;
  ez->speed_stars = (float)sqrt(ez->speed_stars) * STAR_SCALING_FACTOR;

  if (ez->mods & MODS_TOUCH_DEVICE) {
    ez->aim_stars = (float)pow(ez->aim_stars, 0.8f);
  }

  /* calculate total star rating */
  ez->stars = ez->aim_stars + ez->speed_stars +
    (float)fabs(ez->speed_stars - ez->aim_stars) * EXTREME_SCALING_FACTOR;

  return 0;
}

/* taiko diff calc ----------------------------------------------------- */

#define TAIKO_STAR_SCALING_FACTOR 0.04125f
#define TAIKO_TYPE_CHANGE_BONUS 0.75f /* object type change bonus */
#define TAIKO_RHYTHM_CHANGE_BONUS 1.0f
#define TAIKO_RHYTHM_CHANGE_BASE_THRESHOLD 0.2f
#define TAIKO_RHYTHM_CHANGE_BASE 2.0f

typedef struct taiko_object {
  int hit;
  float strain;
  float time;
  float time_elapsed;
  int rim;
  int same_since; /* streak of hits of the same type (rim/center) */
  /*
   * was the last hit type change at an even same_since count?
   * -1 if there is no previous switch (for example if the
   * previous object was not a hit
   */
  int last_switch_even;
} taiko_object_t;

/* object type change bonus */
float taiko_change_bonus(taiko_object_t* cur, taiko_object_t* prev) {
  if (prev->rim != cur->rim) {
    cur->last_switch_even = prev->same_since % 2 == 0;
    if (prev->last_switch_even != cur->last_switch_even &&
        prev->last_switch_even != -1) {
      return TAIKO_TYPE_CHANGE_BONUS;
    }
  } else {
    cur->last_switch_even = prev->last_switch_even;
    cur->same_since = prev->same_since + 1;
  }
  return 0;
}

/* rhythm change bonus */
float taiko_rhythm_bonus(taiko_object_t* cur, taiko_object_t* prev) {
  float ratio;
  float diff;

  if (cur->time_elapsed == 0 || prev->time_elapsed == 0) {
    return 0;
  }

  ratio = al_max(prev->time_elapsed / cur->time_elapsed,
    cur->time_elapsed / prev->time_elapsed);

  if (ratio >= 8) {
    return 0;
  }

  /* this is log base TAIKO_RHYTHM_CHANGE_BASE of ratio */
  diff = (float)fmod(log(ratio) / log(TAIKO_RHYTHM_CHANGE_BASE), 1.0f);

  /*
   * threshold that determines whether the rhythm changed enough
   * to be worthy of the bonus
   */
  if (diff > TAIKO_RHYTHM_CHANGE_BASE_THRESHOLD &&
    diff < 1 - TAIKO_RHYTHM_CHANGE_BASE_THRESHOLD)
  {
    return TAIKO_RHYTHM_CHANGE_BONUS;
  }

  return 0;
}

void taiko_strain(taiko_object_t* cur, taiko_object_t* prev) {
  float decay;
  float addition = 1.0f;
  float factor = 1.0f;

  decay = (float)pow(decay_base[0], cur->time_elapsed / 1000.0f);

  /*
   * we only have strains for hits, also ignore objects that are
   * more than 1 second apart
   */
  if (prev->hit && cur->hit && cur->time - prev->time < 1000.0f) {
    addition += taiko_change_bonus(cur, prev);
    addition += taiko_rhythm_bonus(cur, prev);
  }

  /* 300+bpm streams nerf? */
  if (cur->time_elapsed < 50.0f) {
    factor = 0.4f + 0.6f * cur->time_elapsed / 50.0f;
  }

  cur->strain = prev->strain * decay + addition * factor;
}

void swap_ptrs(void** a, void** b) {
  void* tmp;
  tmp = *a;
  *a = *b;
  *b = tmp;
}

int d_taiko(ezpp_t ez) {
  int i, result;

  /* this way we can swap cur and prev without copying */
  taiko_object_t curprev[2];
  taiko_object_t* cur = &curprev[0];
  taiko_object_t* prev = &curprev[1];

  ez->highest_strains.len = 0;
  ez->max_strain = 0.0f;
  ez->interval_end = STRAIN_STEP * ez->speed_mul;

  /*
   * TODO: separate taiko conversion into its own function
   * so that it can be reused? probably slower, but cleaner,
   * more modular and more readable
   */
  for (i = 0; i < ez->nobjects; ++i) {
    object_t* o = &ez->objects.data[i];

    cur->hit = (o->type & OBJ_CIRCLE) != 0;
    cur->time = o->time;

    if (i > 0) {
      cur->time_elapsed = (cur->time - prev->time) / ez->speed_mul;
    } else {
      cur->time_elapsed = 0;
    }

    if (!o->sound_types) {
      return ERR_SYNTAX;
    }

    cur->strain = 1;
    cur->same_since = 1;
    cur->last_switch_even = -1;
    cur->rim = (o->sound_types[0] & (SOUND_CLAP|SOUND_WHISTLE)) != 0;

    if (ez->original_mode == MODE_TAIKO) {
      goto continue_loop;
    }

    if (o->type & OBJ_SLIDER) {
      /* TODO: too much indentation, pull this out */
      int isound = 0;
      float j;

      /* drum roll, ignore */
      if (!o->slider_is_drum_roll || i == 0) {
        goto continue_loop;
      }

      /*
       * sliders that meet the requirements will
       * become streams of the slider's tick rate
       */
      for (j = o->time;
           j < o->time + o->duration + o->tick_spacing / 8;
           j += o->tick_spacing)
      {
        int sound_type = o->sound_types[isound];
        cur->rim = (sound_type & (SOUND_CLAP | SOUND_WHISTLE)) != 0;
        cur->hit = 1;
        cur->time = j;

        cur->time_elapsed = (cur->time - prev->time) / ez->speed_mul;
        cur->strain = 1;
        cur->same_since = 1;
        cur->last_switch_even = -1;

        /* update strains for this hit */
        if (i > 0 || j > o->time) {
          taiko_strain(cur, prev);
        }

        result = d_update_max_strains(ez, decay_base[0], cur->time,
          prev->time, cur->strain, prev->strain, i == 0 && j == o->time);
        /* warning: j check might fail, floatcheck this */

        if (result < 0) {
          return result;
        }

        /* loop through the slider's sounds */
        ++isound;
        isound %= o->nsound_types;

        swap_ptrs((void**)&prev, (void**)&cur);
      }

      /*
       * since we processed the slider as multiple hits,
       * we must skip the prev/cur swap which we already did
       * in the above loop
       */
      continue;
    }

continue_loop:
    /* update strains for hits and other object types */
    if (i > 0) {
      taiko_strain(cur, prev);
    }

    result = d_update_max_strains(ez, decay_base[0], cur->time, prev->time,
      cur->strain, prev->strain, i == 0);

    if (result < 0) {
      return result;
    }

    swap_ptrs((void**)&prev, (void**)&cur);
  }

  d_weigh_strains(ez, &ez->speed_stars, 0);
  ez->speed_stars *= TAIKO_STAR_SCALING_FACTOR;
  ez->stars = ez->speed_stars;

  return 0;
}

int d_calc(ezpp_t ez) {
  switch (ez->mode) {
    case MODE_STD: return d_std(ez);
    case MODE_TAIKO: return d_taiko(ez);
  }
  info("this gamemode is not yet supported\n");
  return ERR_NOTIMPLEMENTED;
}

/* acc calc ------------------------------------------------------------ */

float acc_calc(int n300, int n100, int n50, int misses) {
  int total_hits = n300 + n100 + n50 + misses;
  float acc = 0;
  if (total_hits > 0) {
    acc = (n50 * 50.0f + n100 * 100.0f + n300 * 300.0f)
      / (total_hits * 300.0f);
  }
  return acc;
}

void acc_round(float acc_percent, int nobjects, int misses, int* n300,
  int* n100, int* n50)
{
  int max300;
  float maxacc;
  misses = al_min(nobjects, misses);
  max300 = nobjects - misses;
  maxacc = acc_calc(max300, 0, 0, misses) * 100.0f;
  acc_percent = al_max(0.0f, al_min(maxacc, acc_percent));
  *n50 = 0;

  /* just some black magic maths from wolfram alpha */
  *n100 = (int)al_round(
    -3.0f * ((acc_percent * 0.01f - 1.0f) * nobjects + misses) * 0.5f
  );

  if (*n100 > nobjects - misses) {
    /* acc lower than all 100s, use 50s */
    *n100 = 0;
    *n50 = (int)al_round(
      -6.0f * ((acc_percent * 0.01f - 1.0f) * nobjects + misses) * 0.2f
    );
    *n50 = al_min(max300, *n50);
  } else {
    *n100 = al_min(max300, *n100);
  }

  *n300 = nobjects - *n100 - *n50 - misses;
}

float taiko_acc_calc(int n300, int n150, int nmiss) {
  int total_hits = n300 + n150 + nmiss;
  float acc = 0;
  if (total_hits > 0) {
    acc = (n150 * 150.0f + n300 * 300.0f) / (total_hits * 300.0f);
  }
  return acc;
}

void taiko_acc_round(float acc_percent, int nobjects, int nmisses,
  int* n300, int* n150)
{
  int max300;
  float maxacc;
  nmisses = al_min(nobjects, nmisses);
  max300 = nobjects - nmisses;
  maxacc = acc_calc(max300, 0, 0, nmisses) * 100.0f;
  acc_percent = al_max(0.0f, al_min(maxacc, acc_percent));
  /* just some black magic maths from wolfram alpha */
  *n150 = (int)al_round(
    -2.0f * ((acc_percent * 0.01f - 1.0f) * nobjects + nmisses)
  );
  *n150 = al_min(max300, *n150);
  *n300 = nobjects - *n150 - nmisses;
}

/* std pp calc --------------------------------------------------------- */

/* some kind of formula to get a base pp value from stars */
float base_pp(float stars) {
  return (float)pow(5.0f * al_max(1.0f, stars / 0.0675f) - 4.0f, 3.0f)
    / 100000.0f;
}

int pp_std(ezpp_t ez) {
  int ncircles = ez->ncircles;
  float nobjects_over_2k = ez->nobjects / 2000.0f;
  float length_bonus = (
    0.95f +
    0.4f * al_min(1.0f, nobjects_over_2k) +
    (ez->nobjects > 2000 ? (float)log10(nobjects_over_2k) * 0.5f : 0.0f)
  );
  float miss_penality = (float)pow(0.97f, ez->nmiss);
  float combo_break = (
    (float)pow(ez->combo, 0.8f) / (float)pow(ez->max_combo, 0.8f)
  );
  float ar_bonus;
  float final_multiplier;
  float acc_bonus, od_bonus;
  float od_squared;
  float hd_bonus;

  /* acc used for pp is different in scorev1 because it ignores sliders */
  float real_acc;
  float accuracy;

  ez->nspinners = ez->nobjects - ez->nsliders - ez->ncircles;

  if (ez->max_combo <= 0) {
    info("W: max_combo <= 0, changing to 1\n");
    ez->max_combo = 1;
  }

  accuracy = acc_calc(ez->n300, ez->n100, ez->n50, ez->nmiss);

  /*
   * scorev1 ignores sliders and spinners since they are free 300s
   * can go negative if we miss everything so we must clamp it
   */

  switch (ez->score_version) {
    case 1:
      real_acc = acc_calc(
        al_max(0, ez->n300 - ez->nsliders - ez->nspinners),
        ez->n100, ez->n50, ez->nmiss);
      break;
    case 2:
      real_acc = accuracy;
      ncircles = ez->nobjects;
      break;
    default:
      info("unsupported scorev%d\n", ez->score_version);
      return ERR_NOTIMPLEMENTED;
  }

  /* ar bonus -------------------------------------------------------- */
  ar_bonus = 1.0f;

  /* high ar bonus */
  if (ez->ar > 10.33f) {
    ar_bonus += 0.3f * (ez->ar - 10.33f);
  }

  /* low ar bonus */
  else if (ez->ar < 8.0f) {
    ar_bonus += 0.01f * (8.0f - ez->ar);
  }

  /* aim pp ---------------------------------------------------------- */
  ez->aim_pp = base_pp(ez->aim_stars);
  ez->aim_pp *= length_bonus;
  ez->aim_pp *= miss_penality;
  ez->aim_pp *= combo_break;
  ez->aim_pp *= ar_bonus;

  /* hidden */
  hd_bonus = 1.0f;
  if (ez->mods & MODS_HD) {
    hd_bonus += 0.04f * (12.0f - ez->ar);
  }

  ez->aim_pp *= hd_bonus;

  /* flashlight */
  if (ez->mods & MODS_FL) {
    float fl_bonus = 1.0f + 0.35f * al_min(1.0f, ez->nobjects / 200.0f);
    if (ez->nobjects > 200) {
      fl_bonus += 0.3f * al_min(1, (ez->nobjects - 200) / 300.0f);
    }
    if (ez->nobjects > 500) {
      fl_bonus += (ez->nobjects - 500) / 1200.0f;
    }
    ez->aim_pp *= fl_bonus;
  }

  /* acc bonus (bad aim can lead to bad acc) */
  acc_bonus = 0.5f + accuracy / 2.0f;

  /* od bonus (high od requires better aim timing to acc) */
  od_squared = (float)pow(ez->od, 2);
  od_bonus = 0.98f + od_squared / 2500.0f;

  ez->aim_pp *= acc_bonus;
  ez->aim_pp *= od_bonus;

  /* speed pp -------------------------------------------------------- */
  ez->speed_pp = base_pp(ez->speed_stars);
  ez->speed_pp *= length_bonus;
  ez->speed_pp *= miss_penality;
  ez->speed_pp *= combo_break;
  if (ez->ar > 10.33f) {
    ez->speed_pp *= ar_bonus;
  }
  ez->speed_pp *= hd_bonus;

  /* scale the speed value with accuracy slightly */
  ez->speed_pp *= 0.02f + accuracy;

  /* it's important to also consider accuracy difficulty when doing that */
  ez->speed_pp *= 0.96f + (od_squared / 1600.0f);

  /* acc pp ---------------------------------------------------------- */
  /* arbitrary values tom crafted out of trial and error */
  ez->acc_pp = (float)pow(1.52163f, ez->od) *
    (float)pow(real_acc, 24.0f) * 2.83f;

  /* length bonus (not the same as speed/aim length bonus) */
  ez->acc_pp *= al_min(1.15f, (float)pow(ncircles / 1000.0f, 0.3f));

  if (ez->mods & MODS_HD) ez->acc_pp *= 1.08f;
  if (ez->mods & MODS_FL) ez->acc_pp *= 1.02f;

  /* total pp -------------------------------------------------------- */
  final_multiplier = 1.12f;
  if (ez->mods & MODS_NF) final_multiplier *= 0.90f;
  if (ez->mods & MODS_SO) final_multiplier *= 0.95f;

  ez->pp = (float)(
    pow(
      pow(ez->aim_pp, 1.1f) +
      pow(ez->speed_pp, 1.1f) +
      pow(ez->acc_pp, 1.1f),
      1.0f / 1.1f
    ) * final_multiplier
  );

  ez->accuracy_percent = accuracy * 100.0f;

  return 0;
}

/* taiko pp calc ------------------------------------------------------- */

int pp_taiko(ezpp_t ez) {
  float length_bonus;
  float final_multiplier;
  float accuracy;

  ez->n300 = al_max(0, ez->max_combo - ez->n100 - ez->nmiss);
  accuracy = taiko_acc_calc(ez->n300, ez->n100, ez->nmiss);

  /* base acc pp */
  ez->acc_pp = (float)pow(150.0f / ez->odms, 1.1f);
  ez->acc_pp *= (float)pow(accuracy, 15.0f) * 22.0f;

  /* length bonus */
  ez->acc_pp *= al_min(1.15f, (float)pow(ez->max_combo / 1500.0f, 0.3f));

  /* base speed pp */
  ez->speed_pp = (
    (float)pow(5.0f * al_max(1.0f, ez->stars / 0.0075f) - 4.0f, 2.0f)
  ) / 100000.0f;

  /* length bonus (not the same as acc length bonus) */
  length_bonus = 1.0f + 0.1f * al_min(1.0f, ez->max_combo / 1500.0f);
  ez->speed_pp *= length_bonus;

  /* miss penality */
  ez->speed_pp *= (float)pow(0.985f, ez->nmiss);

  if (ez->max_combo > 0) {
    ez->speed_pp *= (
      al_min((float)pow(ez->combo - ez->nmiss, 0.5f)
      / (float)pow(ez->max_combo, 0.5f), 1.0f)
    );
  }

  /* speed mod bonuses */
  if (ez->mods & MODS_HD) {
    ez->speed_pp *= 1.025f;
  }

  if (ez->mods & MODS_FL) {
    ez->speed_pp *= 1.05f * length_bonus;
  }

  /* acc scaling */
  ez->speed_pp *= accuracy;

  /* overall multipliers */
  final_multiplier = 1.1f;
  if (ez->mods & MODS_NF) final_multiplier *= 0.90f;
  if (ez->mods & MODS_HD) final_multiplier *= 1.10f;

  ez->pp = (
    (float)pow(
      pow(ez->speed_pp, 1.1f) +
      pow(ez->acc_pp, 1.1f),
      1.0f / 1.1f
    ) * final_multiplier
  );

  ez->accuracy_percent = accuracy * 100.0f;

  return 0;
}

/* main interface ------------------------------------------------------ */

int params_from_map(ezpp_t ez) {
  int res;

  ez->ar = ez->cs = ez->hp = ez->od = 5.0f;
  ez->sv = ez->tick_rate = 1.0f;

  ez->p_flags = 0;
  if (ez->mode_override) {
    ez->p_flags |= P_OVERRIDE_MODE;
  }

  if (ez->data) {
    res = p_map_mem(ez, ez->data, ez->data_size);
  } else if (!strcmp(ez->map, "-")) {
    res = p_map(ez, stdin);
  } else {
    FILE* f = fopen(ez->map, "rb");
    if (!f) {
      perror("fopen");
      res = ERR_IO;
    } else {
      res = p_map(ez, f);
      fclose(f);
    }
  }

  if (res < 0) {
    goto cleanup;
  }

  if (!ez->aim_stars && !ez->speed_stars) {
    res = d_calc(ez);
    if (res < 0) {
      goto cleanup;
    }
  }

cleanup:
  return res;
}

int calc(ezpp_t ez) {
  int res;

  if (!ez->max_combo && (ez->map || ez->data)) {
    if (ez->flags & AUTOCALC_BIT) {
      ez->base_ar = ez->base_od = ez->base_cs = ez->base_hp = -1;
    }
    res = params_from_map(ez);
    if (res < 0) {
      return res;
    }
  } else {
    if (ez->base_ar >= 0) ez->ar = ez->base_ar;
    if (ez->base_od >= 0) ez->od = ez->base_od;
    if (ez->base_cs >= 0) ez->cs = ez->base_cs;
    if (ez->base_hp >= 0) ez->hp = ez->base_hp;
    mods_apply(ez);
  }

  if (ez->mode == MODE_TAIKO) {
    ez->stars = ez->speed_stars;
  }

  if (ez->accuracy_percent >= 0) {
    switch (ez->mode) {
      case MODE_STD:
        acc_round(ez->accuracy_percent, ez->nobjects, ez->nmiss,
          &ez->n300, &ez->n100, &ez->n50);
        break;
      case MODE_TAIKO:
        taiko_acc_round(ez->accuracy_percent, ez->max_combo,
          ez->nmiss, &ez->n300, &ez->n100);
        break;
    }
  }

  if (ez->combo < 0) {
    ez->combo = ez->max_combo - ez->nmiss;
  }

  ez->n300 = ez->nobjects - ez->n100 - ez->n50 - ez->nmiss;

  switch (ez->mode) {
    case MODE_STD: res = pp_std(ez); break;
    case MODE_TAIKO: res = pp_taiko(ez); break;
    default:
      info("pp calc for this mode is not yet supported\n");
      return ERR_NOTIMPLEMENTED;
  }

  if (res < 0) {
    return res;
  }

  return 0;
}

OPPAIAPI
ezpp_t ezpp_new(void) {
  ezpp_t ez = calloc(sizeof(struct ezpp), 1);
  if (ez) {
    ez->mode = MODE_STD;
    ez->mods = MODS_NOMOD;
    ez->combo = -1;
    ez->score_version = 1;
    ez->accuracy_percent = -1;
    ez->base_ar = ez->base_od = ez->base_cs = ez->base_hp = -1;
    array_reserve(&ez->objects, 600);
    array_reserve(&ez->timing_points, 16);
    array_reserve(&ez->highest_strains, 600);
  }
  return ez;
}

void free_owned_map(ezpp_t ez) {
  if (ez->flags & OWNS_MAP_BIT) {
    free(ez->map);
    free(ez->data);
    ez->flags &= ~OWNS_MAP_BIT;
  }
  ez->map = 0;
  ez->data = 0;
  ez->data_size = 0;
  if (ez->flags & AUTOCALC_BIT) {
    ez->max_combo = 0; /* force re-parse */
  }
}

OPPAIAPI
void ezpp_free(ezpp_t ez) {
  free_owned_map(ez);
  array_free(&ez->objects);
  array_free(&ez->timing_points);
  array_free(&ez->highest_strains);
  m_free(ez);
  free(ez);
}

OPPAIAPI
int ezpp(ezpp_t ez, char* mapfile) {
  free_owned_map(ez);
  ez->map = mapfile;
  return calc(ez);
}

OPPAIAPI
int ezpp_data(ezpp_t ez, char* data, int data_size) {
  free_owned_map(ez);
  ez->data = data;
  ez->data_size = data_size;
  return calc(ez);
}

void* memclone(void* p, int size) {
  void* res = malloc(size);
  if (res) memcpy(res, p, size);
  return res;
}

char* strclone(char* s) {
  int len = (int)strlen(s) + 1;
  return memclone(s, len);
}

OPPAIAPI
int ezpp_dup(ezpp_t ez, char* mapfile) {
  free_owned_map(ez);
  ez->flags |= OWNS_MAP_BIT;
  ez->map = strclone(mapfile);
  return calc(ez);
}

OPPAIAPI
int ezpp_data_dup(ezpp_t ez, char* data, int data_size) {
  free_owned_map(ez);
  ez->flags |= OWNS_MAP_BIT;
  ez->data = memclone(data, data_size);
  ez->data_size = data_size;
  return calc(ez);
}

OPPAIAPI float ezpp_pp(ezpp_t ez) { return ez->pp; }
OPPAIAPI float ezpp_stars(ezpp_t ez) { return ez->stars; }
OPPAIAPI int ezpp_mode(ezpp_t ez) { return ez->mode; }
OPPAIAPI int ezpp_combo(ezpp_t ez) { return ez->combo; }
OPPAIAPI int ezpp_max_combo(ezpp_t ez) { return ez->max_combo; }
OPPAIAPI int ezpp_mods(ezpp_t ez) { return ez->mods; }
OPPAIAPI int ezpp_score_version(ezpp_t ez) { return ez->score_version; }
OPPAIAPI float ezpp_aim_stars(ezpp_t ez) { return ez->aim_stars; }
OPPAIAPI float ezpp_speed_stars(ezpp_t ez) { return ez->speed_stars; }
OPPAIAPI float ezpp_aim_pp(ezpp_t ez) { return ez->aim_pp; }
OPPAIAPI float ezpp_speed_pp(ezpp_t ez) { return ez->speed_pp; }
OPPAIAPI float ezpp_acc_pp(ezpp_t ez) { return ez->acc_pp; }
OPPAIAPI float ezpp_accuracy_percent(ezpp_t ez) { return ez->accuracy_percent; }
OPPAIAPI int ezpp_n300(ezpp_t ez) { return ez->n300; }
OPPAIAPI int ezpp_n100(ezpp_t ez) { return ez->n100; }
OPPAIAPI int ezpp_n50(ezpp_t ez) { return ez->n50; }
OPPAIAPI int ezpp_nmiss(ezpp_t ez) { return ez->nmiss; }
OPPAIAPI char* ezpp_title(ezpp_t ez) { return ez->title; }
OPPAIAPI char* ezpp_title_unicode(ezpp_t ez) { return ez->title_unicode; }
OPPAIAPI char* ezpp_artist(ezpp_t ez) { return ez->artist; }
OPPAIAPI char* ezpp_artist_unicode(ezpp_t ez) { return ez->artist_unicode; }
OPPAIAPI char* ezpp_creator(ezpp_t ez) { return ez->creator; }
OPPAIAPI char* ezpp_version(ezpp_t ez) { return ez->version; }
OPPAIAPI int ezpp_ncircles(ezpp_t ez) { return ez->ncircles; }
OPPAIAPI int ezpp_nsliders(ezpp_t ez) { return ez->nsliders; }
OPPAIAPI int ezpp_nspinners(ezpp_t ez) { return ez->nspinners; }
OPPAIAPI int ezpp_nobjects(ezpp_t ez) { return ez->nobjects; }
OPPAIAPI float ezpp_ar(ezpp_t ez) { return ez->ar; }
OPPAIAPI float ezpp_cs(ezpp_t ez) { return ez->cs; }
OPPAIAPI float ezpp_od(ezpp_t ez) { return ez->od; }
OPPAIAPI float ezpp_hp(ezpp_t ez) { return ez->hp; }
OPPAIAPI float ezpp_odms(ezpp_t ez) { return ez->odms; }
OPPAIAPI int ezpp_autocalc(ezpp_t ez) { return ez->flags & AUTOCALC_BIT; }

OPPAIAPI float ezpp_time_at(ezpp_t ez, int i) {
  return ez->objects.len ? ez->objects.data[i].time : 0;
}

OPPAIAPI float ezpp_strain_at(ezpp_t ez, int i, int difficulty_type) {
  return ez->objects.len ? ez->objects.data[i].strains[difficulty_type] : 0;
}

#define setter(t, x) \
OPPAIAPI void ezpp_set_##x(ezpp_t ez, t x) { \
  ez->x = x; \
  if (ez->flags & AUTOCALC_BIT) { \
    calc(ez); \
  } \
}
setter(float, aim_stars)
setter(float, speed_stars)
setter(float, base_ar)
setter(float, base_od)
setter(float, base_hp)
setter(int, mode)
setter(int, combo)
setter(int, score_version)
setter(float, accuracy_percent)
#undef setter

OPPAIAPI
void ezpp_set_autocalc(ezpp_t ez, int autocalc) {
  if (autocalc) {
    ez->flags |= AUTOCALC_BIT;
  } else {
    ez->flags &= ~AUTOCALC_BIT;
  }
}

OPPAIAPI
void ezpp_set_mods(ezpp_t ez, int mods) {
  if ((mods ^ ez->mods) & (MODS_MAP_CHANGING | MODS_SPEED_CHANGING)) {
    /* force map reparse */
    ez->aim_stars = ez->speed_stars = ez->stars = 0;
    ez->max_combo = 0;
  }
  ez->mods = mods;
  if (ez->flags & AUTOCALC_BIT) {
    calc(ez);
  }
}

#define clobber_setter(t, x) \
OPPAIAPI \
void ezpp_set_##x(ezpp_t ez, t x) { \
  ez->aim_stars = ez->speed_stars = ez->stars = 0; \
  ez->max_combo = 0; \
  ez->x = x; \
  if (ez->flags & AUTOCALC_BIT) { \
    calc(ez); \
  } \
}
clobber_setter(float, base_cs)
clobber_setter(int, mode_override)
#undef clobber_setter

#define acc_clobber_setter(t, x) \
OPPAIAPI \
void ezpp_set_##x(ezpp_t ez, t x) { \
  ez->accuracy_percent = -1; \
  ez->aim_stars = ez->speed_stars = ez->stars = 0; \
  ez->max_combo = 0; \
  ez->x = x; \
  if (ez->flags & AUTOCALC_BIT) { \
    calc(ez); \
  } \
}

acc_clobber_setter(int, nmiss)
acc_clobber_setter(int, end)
acc_clobber_setter(float, end_time)

OPPAIAPI
void ezpp_set_accuracy(ezpp_t ez, int n100, int n50) {
  ez->accuracy_percent = -1;
  ez->n100 = n100;
  ez->n50 = n50;
  if (ez->flags & AUTOCALC_BIT) {
    calc(ez);
  }
}

#endif /* OPPAI_IMPLEMENTATION */
