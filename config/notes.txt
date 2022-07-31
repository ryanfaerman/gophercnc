How does config get set on a general level?

Something like git? --global user.email

Something like:

table: config
keys:
  - scope
  - key
  - value

log.level = info
machine.active = mpcnc
machine.geometry.x_axis = 200
machine.geometry.y_axis = 300

How would limits and loading the current stuff work for toolpaths?

gophercnc generate operation facing // for example... not sure the UX yet

let's say a facing is requested larger than the machine can handle.

requested_x = 200
requested_y = 3000

it shoudl fail... how would that work in code:

if requested_y > config.Get("machine.geometry.x_axis") {}

That seems workable... but seems error prone. Ideally I want to have it be typed.

machine := config.MachineConfig()
if machine.ValidateLimits(requested_x, requested_y) {}

This is more promising. In this case, we go through the Config to get the current
machine.Machine instance, where that is defiend as something like:

type Machine struct {
  Geometry struct {
  ...
  }
}

func (m Machine) ValidLimits(x, y) bool {
 return m.Geometry.X <= x && m.Geometry.Y <= y
}


Machines are defineable within Fusion and can be exported as XML. This contains
all the limits and such. But, I'd want to support defining a machine manually,
this way if you don't have everything in fusion you don't need it.

Same with tools, even though right now it is all defined in terms of F360.

Both tools and machines are defined in files. I'd like to not handle truly
importing them, rather, pulling them in on-demand. That means "import" just
copies the contents into a known location and adds it to a database. If the
file is already imported, it should ignore it. Or rather, offer to replace it.

the database should be slightly more generic, so there is one table for all the
file lookups and the like, rather than a new table for every thing that works
in this manner.

Maybe something like:

RESOURCE
kind (tool, machine, etc.)
format (f360.tools, f360.machine, gophercnc.yml, etc.)
name (mpcnc, example-library)
path (machines/banana.machine, tools/mpcnc.tools, etc.)

CONFIG
uri (machine.active, tools.library.active, tools.active_tool)
data (mpcnc, mpcnc.tools, , )


Where for the key "machine.active", we can look it up by going to the
CONFIG.uri(machine.active), then returning "mpcnc". From there, we can go to
the resource table. SELECT path FROM resources WHERE name = 'mpcnc'. With the
path, we can now load the actual file: library.Load(path)

This would go: Config -> Resource -> FILESYSTEM.

Although... with this approach, we can probably just use a JOIN.

 select resources.name, resources.path from config join resources on config.data = resources.name where config.uri = "library.active";

What about for config that isn't a path or reference? Does it even matter?

The config approach with a `uri` and `data` approach... isn't specific. It's basically a key/value pairing. If I wanted to have something like a hook for gcode added at the beginning...

uri: gcode.preamble
data: G1.....

Is this a touch over-engineered? Maybe. But it is fun and not so insane. Plenty of things store config in a database rather than a flat file.

The config can be stored in a flat file and validated against the DB on startup. Maybe no the whole thing... maybe we just store the hash of the file in the database itself. IF it changes... the keys are updated from the config file. If not, we just proceed.

This way it becomes self-updating. And we can dump the config to a similar flat file.
urn:auth0
A TOML file is syntactically similar to a git config file. And it would support
sections, keys, and the types we need. Plus it has well defined parser.

A TOML file can also be the native way that libraries, machines, and their like are defined.

[[tool]]

