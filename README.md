WADmake
=======

A utility for creating WAD, ZIP (PK3) and 7Z (PK7) archives for Doom.

FAQ
---

Q: Why?

A: Because deutex needed to be put out of its misery.

Q: Why not SLADE/XWE?

A: SLADE operates on a single monolithic file, which cannot be version-controlled in a useful fashion.  Revision control has been standard practice in the software field for at least a decade, and the benefits of using revision control software are too numerous to list here.

Q: Why not use zip/7za in a shell script or makefile?

A: WADmake is designed to work with Doom-specific file formats natively.

Q: Why not several dozen single-use utilities?

A: Because I wanted the utility to be a single executable that you could copy around freely and drop into a project, like deutex.  Of course, there are still some operations that require third party utilities, but such utilities are optional if you do not require their functionality.

Q: Why Go?

A: Again, because I wanted the utility to be a single executable.  Go is incredibly good at producing a single statically linked executable.  A prior version of WADmake was written in C++, but I switched since Go is so much easier to maintain.

Q: Why Lua?

A: A build system should be extensible, and I did not want to invent a new extension language from scratch.  In addition to its small size and ubiquity as an extension language, Lua seemed like a natrual fit as it originated as a configuration language and has been used in other build systems like Premake.

License
-------

Currently [GNU AGPLv3](https://www.gnu.org/licenses/agpl-3.0.html).  If there is a compelling case to be made for a more permissive license, I am open to suggestions.
