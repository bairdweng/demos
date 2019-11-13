import { NestFactory } from '@nestjs/core';
import { ApplicationModule } from './app.modules';
import { createConnection } from "typeorm";



createConnection().then(async connection => {
  console.log("Post has been saved: ");
  bootstrap();
}).catch(error => console.log("Error: ", error));



async function bootstrap() {
  const app = await NestFactory.create(ApplicationModule);
  await app.listen(3000);
}




