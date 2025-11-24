<?php

namespace Go\Contract;

/**
 * Defines the contract for all executable PHP jobs.
 */
interface JobInterface
{
    /**
     * Executes the job.
     * @param mixed $payload The data required for the job to run.
     * @return mixed The result of the job.
     */
    public function handle($payload);
}
