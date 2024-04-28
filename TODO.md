# TODO

- [x] Create env first pass.
- [x] Enter env first pass.
- [x] Allow mounting of local dirs to env using docker volume.
  - [ ] Mount home dir handling seamlessly.
- [ ] Handle env management via per project pkl files (need to store container IDs).
- [ ] Mount .ssh configs for seamless git access.
- [ ] Agnostos RC for default env setup.
  - [ ] List of mounted dirs.
  - [ ] List of preferred language versions, defaulting to latest.

# Considerations

- Nix OS could be a good option for repeatable envs if ubuntu is too painful. [read more](https://nixos.org/)

- Fuse allows us to mount remote (container) dir to local dir (i.e an interpreter/compiler).
  - This means we could skip all out container setup and just mount the container dir containing the interpreter/compiler. [read more](https://github.com/libfuse/sshfs)
  - NFS is an alternative option. [read more](https://gist.github.com/proudlygeek/5721498)
    - Faster but lacks encryption.
