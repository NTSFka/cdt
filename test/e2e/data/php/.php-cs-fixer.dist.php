<?php

$finder = PhpCsFixer\Finder::create()
    ->ignoreUnreadableDirs()
    ->exclude('db')
    ->exclude('log')
    ->exclude('temp')
    ->in(__DIR__)
;

return (new PhpCsFixer\Config())
    ->setRules(['@PER-CS2.0' => true])
    ->setFinder($finder)
;
