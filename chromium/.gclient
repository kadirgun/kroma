solutions = [
  {
    "name"        : 'src',
    "url"         : 'https://chromium.googlesource.com/chromium/src.git@refs/tags/147.0.7727.56',
    "deps_file"   : 'DEPS',
    "managed"     : False,
    "custom_deps" : {
    },
    "custom_vars": {
      "checkout_pgo_profiles": True,
    },
  },
]
target_os = ['win']