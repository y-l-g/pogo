<?php

namespace Illuminate\Contracts\Console {
    interface Kernel
    {
        public function bootstrap();
    }
}

namespace Illuminate\Foundation {
    class Application
    {
        private $bindings = [];

        public function make($abstract)
        {
            if (isset($this->bindings[$abstract])) {
                $concrete = $this->bindings[$abstract];
                return $concrete instanceof \Closure ? $concrete($this) : $concrete;
            }
            return new $abstract();
        }

        public function bound($abstract)
        {
            return isset($this->bindings[$abstract]) || class_exists($abstract);
        }

        public function environment()
        {
            return "testing";
        }

        // Minimal implementation of 'call' for Method Injection
        public function call($callback, array $parameters = [], $defaultMethod = null)
        {
            if (is_array($callback)) {
                $method = new \ReflectionMethod($callback[0], $callback[1]);
            } else {
                $method = new \ReflectionFunction($callback);
            }

            $args = [];
            foreach ($method->getParameters() as $param) {
                $name = $param->getName();
                if (array_key_exists($name, $parameters)) {
                    $args[] = $parameters[$name];
                } elseif ($param->getType() && !$param->getType()->isBuiltin()) {
                    $cls = $param->getType()->getName();
                    $args[] = $this->make($cls);
                } else {
                    $args[] = null;
                }
            }

            return call_user_func_array($callback, $args);
        }

        public function bind($abstract, $concrete)
        {
            $this->bindings[$abstract] = $concrete;
        }
    }
}
