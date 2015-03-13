acidforth
=========

A modular synthesizer and sequencer.

Introduction
------------

acidforth is a digital, modular synthesizer. It has 8 independent phase
generator oscillators, 8 envelopes and tons of other modules. In a basic
configuration it can be thought of as a very flexible FM synthesizer.

It also comes with a sequencer. It mimics the functionality of the TB-303
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
    + 2 /

`440` is put on the stack. `op1` pops it from the stack and sets its phase
increment accordingly, and then drops its current phase state on the stack.
`sintab` pops the phase state and uses it to index a sine waveform table,
putting the value for that index on the stack. This is repeated for `op2`,
leaving two values on the stack. `+` adds the two values together, and then it
is all divided by 2, leaving a single value on the stack.

After the program has run, the value that is on top of the stack is output as
an audio sample.

Usage
-----

When acidforth is started, it waits for a program on stdin. When it has a
program it is compiled and run immediately. Example use:

    ./acidforth < patches/cool

Words
-----

### Plumbing

    drop ( x -- )
    dup    ( x -- x x )
    swap   ( x y -- y x)
    rot    ( x y z -- y z x )
    *      ( x y -- x * y )
    +      ( x y -- x + y )
    -      ( x y -- x - y )
    /      ( x y -- x / y )
    %      ( x y -- x % y )
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


Macros
------

Programs can use macros to avoid duplicating code or getting lost in the stack.
These macros look like Forth word definitions but are inlined at compilation.

Example

    : double 2 * ;
    440 op1 double 1 -

Comments
--------

Comments in the source code start with "(" followed by white space and then
the rest of the comment, ending with the character ")", whitespace or not. This
should be familiar if you have ever commented forth source code.
