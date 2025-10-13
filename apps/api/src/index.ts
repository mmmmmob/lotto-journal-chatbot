import { prisma } from '@lotto/db';
import { Elysia } from 'elysia';

const app = new Elysia()
  .get('/', () => 'API running ✅')
  .get('/tickets', async () => {
    const tickets = await prisma.ticket.findMany({
      take: 20,
      orderBy: { createdAt: 'desc' },
    });
    return tickets;
  })
  .listen(8787);

console.log('🧠 API listening on http://localhost:8787');

export default app.fetch;
