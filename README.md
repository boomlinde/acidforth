acidforth
=========

A modular synthesizer and sequencer. Hard to learn, impossible to master.

Introduction
------------

acidforth is a modular synthesizer. It has 8 independent phase generator
oscillators, 8 envelopes and tons of other modules. In a basic configuration
it can be thought of as a very flexible FM synthesizer.

It also comes with a sequencer that mimics the functionality of the TB-303
sequencer, including note slides and accents. The sequencer isn't programmed
by the user in an ordinary fashion. Instead, the user feeds it with a random
seed that determines the pattern structure.

The connections between the modules and sequencer are managed using a simple
RPN stack-based programming language. The programs are executed once per output
sample.

For example, a program to output a sine at 440 Hz mixed with another at 445
might look like this:

    440 op1 sintab
    445 op2 sintab
    + 2 / >out

`440` is put on the stack. `op1` pops it from the stack and sets its phase
increment accordingly, and then drops its current phase state on the stack.
`sintab` pops the phase state and uses it to index a sine waveform table,
putting the value for that index on the stack. This is repeated for `op2`,
leaving two values on the stack. `+` adds the two values together, and then it
is all divided by 2, leaving a single value on the stack. `>out` pops that
value and outputs it to both audio channels.

Usage
-----

When acidforth is started, it loads the program specified in the last command
line parameter, and compiles and runs it immediately. Example use:

    ./acidforth patches/cool

acidforth can also load wave samples. These should be listed as parameters for
the command, e.g.

    ./acidforth samples/*.wav patches/cool

The sequencer can be started and stopped by entering empty lines into stdin.

The option `-l` will list available MIDI interfaces and exit immediately, and
the `-m N` option will let you select one of the interfaces for control input.

acidforth may start an HTTP server by using the option `-s` followed by a
`interface:port` pair.

Words
-----

### Plumbing

    drop   ( x -- )
    dup    ( x -- x x )
    swap   ( x y -- y x)
    rot    ( x y z -- y z x )
    *      ( x y -- x * y )
    +      ( x y -- x + y )
    -      ( x y -- x - y )
    /      ( x y -- x / y )
    %      ( x y -- x % y )
    _      ( x -- floor.x ) 
    clip   ( x -- x clipped to range -1, 1)
    pi     ( -- 3.14... )
    =      ( x y -- 1 if x = y else 0 )
    <      ( x y -- 1 if x < y else 0 )
    >      ( x y -- 1 if x > y else 0 )
    <=     ( x y -- 1 if x <= y else 0 )
    >=     ( x y -- 1 if x >= y else 0 )
    ~      ( x -- ~x )
    .      ( x -- print... )
    push   ( x -- push x to secondary stack )
    pop    ( pop x from secondary stack -- x )
    srate  ( pushes sample rate to stack )
    m2f    ( pops a midi note number and pushes the corresponding frequency )
    sintab ( pops an index value 0 - 1, pushes corresponding sine table value )
    >out   ( pops a value and uses that as the next sample for both channels )
    >out1  ( pops a value and uses that as the next sample for the first channel )
    >out2  ( pops a value and uses that as the next sample for the second channel )
    prompt ( pushes the last entered number on stdin or 0 if none entered )

### Phase generators

    op1 ... op8               ( pops frequency and pushes its current phase )
    op1.rst ... op8.rst       ( pops value  and resets the phase if = 1 )
    op1.cycle? ... op8.cycle? ( push 1 if the op just cycled, otherwise 0 )

### Envelopes

    env1 ... env8     ( pops gate and pushes current envelope output value )
    env1.a ... env8.a ( pops envelope attack length )
    env1.d ... env8.d ( pops envelope decay length )
    env1.r ... env8 r ( pops envelope release length )

### Accumulators

    >mix1 ... >mix4 ( pops and adds up to the accumulator )
    mix1> ... mix4> ( outputs the sum of accumulated values and clears )

### Registers

    >A ... >Z ( pops and stores in temporary variable )
    A> ... Z> ( loads variable and push to stack )

### Sequencer

    seq.pitch   ( push current midi note number )
    seq.gate    ( push current gate state )
    seq.accent  ( push current accent state )
    seq.tune    ( pop and set tune offset from middle c )
    seq.tempo   ( pop and set tempo )
    seq.pattern ( pop and set pattern )
    seq.len     ( pop sequence length )
    seq.trig    ( push sequencer sync pulse to stack)
	seq.swing   ( pop a swing factor for the sequener )

### Drum pattern sequencers

    dseq1 ... dseq8         ( pop pattern from stack, output drum trigger )
    dseq1.len ... dseq8.len ( pop pattern length from stack)

### Discrete value sequencers
	vseq1 ... vseq8         ( pop pattern from stack, output value 0 or 1 )
    vseq1.len ... vseq8.len ( pop pattern length from stack)

### MIDI data

	cc  ( pop n from stack, output last value of MIDI CC n )
	key ( pop n from stack, push 1 or 0, alternating each time note 1 is pressed )
	mom ( pop n from stack, push 1 if note n is held, else 0 )
	vol ( pop n from stack, output last velocity of note n )

### Samples

Samples are loaded into the synthesizer from the positional command line
arguments that precede the synth program. They are each assigned words
according to their base filenames. Loading e.g. `bd.wav` will create the words
`bd.wav` and `bd.wav.rate`.

    samplename ( pops a trigger that resets sample on rising edge and pushes sample value )
    samplename.rate ( multiplies sample speed by top of stack )

### Discrete value and drum sequencers

Discrete value and drum sequencers step through the bits of the pattern value
in time with the main sequencer, one bit per 16th note. When they have stepped
.len steps, they loop.

The sequencers can be multiplexed. That is, a single dseq/vseq word may be
used several times throughout a program. The .len may not be canged multiple
times per program.

The difference between the two is that the drum sequencer will only output
its step value for a 32nd, after which it will output 0, while the discrete
value sequencer will output its value for the full 16th.

Macros
------

Programs can use macros to avoid duplicating code or getting lost in the stack.
These macros look like Forth word definitions but are inlined at compilation.

Example

    : double 2 * ;
    440 op1 double 1 - >out

Comments
--------

Comments in the source code start with "(" followed by white space and then
the rest of the comment, ending with the character ")", whitespace or not. This
should be familiar if you have ever commented forth source code.

MIDI control
------------

* `cc`: continuous controller #n. The output value is in the range 0-1.
* `key` toggle key switch. The output value toggles between 0 and 1 every
  time a note on event for the note `n` is received.
* `mom` momentary key switch. The output value is 1 while the note `n` is on
  and 0 when it is off.
* `vel` note velocity. The output value is in the range 0-1 and corresponds to
  the velocity value of the last press of note `n`.

HTTP control
------------

The HTTP server can be used to toggle sequencer playback and to update the
running program.

* `/compiler` (POST) will compile the program in the body of the request
  and respond with a 200 status if successful.
* `/playback` (POST) will toggle the playback of the sequencer and respond
  with a 200 status. The body of the request will be ignored.

Building
--------

acidforth depends on "github.com/rakyll/portmidi" and
"github.com/gordonklaus/portaudio". These require portmidi (version 217) and
portaudio (v19) and their headers, along with a C compiler for use with cgo.
The portmidi package additionally uses pkg-config to find the library to link
to. All these instructions assume that you have Go installed.

Install the external dependencies. After that, and setting up `GOPATH` you can
use `go get github.com/boomlinde/acidforth` to build the acidforth binary.

### Void Linux

    # xbps-install gcc portaudio-devel portmidi-devel alsa-lib-devel pkg-config

### Ubuntu

I have not tried this but it seems straight forward enough

    # apt-get install pkg-config gcc portaudio19-dev libportmidi-dev

### OS X

Off the top of my head, using the brew package manager...

    $ brew install portaudio portmidi

### Windows

I had a terrible experience trying to build the dependencies, but using TDM-GCC
I could build the portmidi and portaudio dlls and link them to the go wrappers
after some minor modifications of their cgo directives.
 

Box art
-------

![box art](http://i.imgur.com/ODgoorr.png)

Box art courtesy of [seece](https://github.com/seece)!
