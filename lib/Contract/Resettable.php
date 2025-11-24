<?php

namespace Go\Contract;

interface Resettable
{
    /**
     * Reset the job state after execution.
     * This is called by the worker loop before garbage collection.
     * Use this to clear large arrays, close file handles, or reset static properties.
     */
    public function reset(): void;
}
