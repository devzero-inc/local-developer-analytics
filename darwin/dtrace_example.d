#!/usr/sbin/dtrace -s

dtrace:::BEGIN {
    printf("Starting command execution timing\n");
}

proc:::exec-success {
    self->start_time = timestamp;
    self->cmd = execname;
}

proc:::exit {
    /* Check if we have a start time for this process */
    if (self->start_time) {
        this->duration = (timestamp - self->start_time) / 1000000; /* Convert to milliseconds */
        printf("Command '%s' executed in %d ms\n", self->cmd, this->duration);
    }
}

