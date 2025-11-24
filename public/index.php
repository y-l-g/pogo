<?php

require __DIR__ . '/../vendor/autoload.php';

// Si FrankenPHP a déjà lancé le worker (ce qui est le cas ici),
// cette fonction s'attachera au pool existant ou ne fera rien.
Go\start_worker_pool('/app/worker.php');

echo "Envoi du job...<br>";
// ... reste du code
