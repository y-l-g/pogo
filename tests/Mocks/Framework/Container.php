<?php

namespace Framework;

class Container
{
    private static $instance;
    private $services = [];

    public static function getInstance()
    {
        if (self::$instance === null) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    public function bind($name, $instance)
    {
        $this->services[$name] = $instance;
    }

    public function get($name)
    {
        if (!isset($this->services[$name])) {
            throw new \Exception("Service '$name' not found in Container");
        }
        return $this->services[$name];
    }

    public function make($className)
    {
        if (!class_exists($className)) {
            throw new \Exception("Class $className not found");
        }
        if (method_exists($className, 'create')) {
            return $className::create($this);
        }
        return new $className();
    }

    /**
     * Invokes a callable, resolving dependencies from the container or parameters.
     *
     * @param callable|array $callback
     * @param array $parameters Named parameters from the JSON payload
     */
    public function call($callback, array $parameters = [])
    {
        if (is_array($callback)) {
            $reflector = new \ReflectionMethod($callback[0], $callback[1]);
        } else {
            $reflector = new \ReflectionFunction($callback);
        }

        $dependencies = [];
        foreach ($reflector->getParameters() as $param) {
            $name = $param->getName();
            $type = $param->getType();

            if (array_key_exists($name, $parameters)) {
                // 1. Explicit parameter provided in payload
                $dependencies[] = $parameters[$name];
            } elseif ($type && !$type->isBuiltin()) {
                // 2. Dependency Injection (Service)
                $serviceName = $type->getName();
                // Simple mapping for our mock container: 'Framework\Logger' -> 'logger'
                if ($serviceName === 'Framework\Logger') {
                    $dependencies[] = $this->get('logger');
                } else {
                    $dependencies[] = $this->make($serviceName);
                }
            } elseif ($param->isDefaultValueAvailable()) {
                // 3. Default value
                $dependencies[] = $param->getDefaultValue();
            } else {
                throw new \Exception("Unable to resolve parameter '{$name}'");
            }
        }

        return call_user_func_array($callback, $dependencies);
    }
}
