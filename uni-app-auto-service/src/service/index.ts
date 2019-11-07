import { NestFactory } from '@nestjs/core';
import * as express from 'express';
import { ApplicationModule } from './app.modules';

const server = express();

const app = NestFactory.create(ApplicationModule, server);
app.listen(3000, () => {
  console.log('Typescript Nest app & Express server running on port 3000');
});