// scriptish is a library to help you port bash scripts to Golang
//
// inspired by:
//
// - http://labix.org/pipe
// - https://github.com/bitfield/script
//
// Copyright 2019-present Ganbaro Digital Ltd
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions
// are met:
//
//   * Redistributions of source code must retain the above copyright
//     notice, this list of conditions and the following disclaimer.
//
//   * Redistributions in binary form must reproduce the above copyright
//     notice, this list of conditions and the following disclaimer in
//     the documentation and/or other materials provided with the
//     distribution.
//
//   * Neither the names of the copyright holders nor the names of his
//     contributors may be used to endorse or promote products derived
//     from this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS
// FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE
// COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT,
// INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING,
// BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
// LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT
// LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN
// ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
// POSSIBILITY OF SUCH DAMAGE.

package scriptish

// StepOption represents a functional option that can be
// used to reconfigure a Pipe before we execute a PipeCommand.
type StepOption struct {
	// we will call this before pipe.RunCommand() is called
	//
	// use this to do any setup work (eg output redirection)
	runSetup Command

	// we will call this after pipe.RunCommand() has been called
	//
	// use this to clean up afterwards (eg close any files that
	// were opened during the setup phase)
	runTeardown Command
}

// noopCommand is a dummy Command, that we use to avoid having to check
// for nil StepOption commands.
func noopCommand(p *Pipe) (int, error) {
	return StatusOkay, nil
}

// NewStepOption wraps up your setup and teardown phases for you.
func NewStepOption(setupPhase Command, teardownPhase Command) *StepOption {
	// make sure we have the phases
	if setupPhase == nil {
		setupPhase = noopCommand
	}
	if teardownPhase == nil {
		teardownPhase = noopCommand
	}

	return &StepOption{
		runSetup:    setupPhase,
		runTeardown: teardownPhase,
	}
}

// ApplySetupPhasesToPipe executes the setup phases of the given StepOptions.
//
// If any of the setup phases fail, ApplySetupPhasesToPipe returns an error,
// and does not execute any remaining setup phases.
func ApplySetupPhasesToPipe(p *Pipe, opts ...*StepOption) (int, error) {
	for _, opt := range opts {
		// apply the option
		p.RunCommand(opt.runSetup)

		// we stop executing the moment something goes wrong
		err := p.Error()
		if err != nil {
			// debugging support
			statusCode := p.StatusCode()
			Tracef("status code: %d", statusCode)
			Tracef("error: %s", err.Error())

			// we cannot continue
			return statusCode, err
		}
	}

	return StatusOkay, nil
}

// ApplyTeardownPhasesToPipe executes the teardown phases of the given
// StepOptions, in reverse order.
func ApplyTeardownPhasesToPipe(p *Pipe, opts ...*StepOption) {
	for i := len(opts) - 1; i >= 0; i-- {
		// apply the option
		//
		// NOTE that we do not use pipe.RunCommand() here, because we
		// do not want the teardown phase to interfere with the error
		// status of the pipe!
		opts[i].runTeardown(p)

		// we don't stop stop executing the moment something goes wrong,
		// because we want the other teardown work to at least try
		// to happen
	}
}
