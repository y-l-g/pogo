<?php

require __DIR__ . '/../vendor/autoload.php';

require_once __DIR__ . '/HelloWorldJob.php';

use Go\Runtime\Protocol;

(new Protocol())->run();