#!/usr/bin/env node
import { Command } from 'commander';
import packageJson from '../package.json';
import { helloWorld } from './commands/helloWorld';

const program = new Command();

program
  .name('hello-world')
  .description('A starter template for npm CLI applications')
  .version(packageJson.version, '-v, --version', 'output the version number');

program
  .command('hello-world')
  .description('Print a greeting')
  .action(() => {
    console.log(helloWorld());
  });

program.parse();
