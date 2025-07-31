import { Hono } from "hono";
import { getCelery } from "./celery";

const app = new Hono();

app.get("/", (c) => {
	return c.text("Hello, World!");
});

app.get("/cel/example", async (c) => {
	const start = performance.now();
	try {
		const celery = await getCelery();
		const you = c.req.query("you") ?? "stranger";
		const result = celery.eval("i.greet(you)", { i: "CEL", you });
		console.log(`CEL evaluation took ${performance.now() - start}ms`);
		return c.json(result);
	} catch (error) {
		return c.json({ error: `CEL evaluation failed: ${error}` }, 500);
	}
});

app.get("/cel/overdraft-test", async (c) => {
	try {
		const celery = await getCelery();

		const balance = Number(c.req.query("balance")) || 1;
		const overdraftProtection = c.req.query("overdraftProtection") === "true";
		const overdraftLimit = Number(c.req.query("overdraftLimit")) || 10000;
		const withdrawal = Number(c.req.query("withdrawal")) || 5000;

		const data = celery.eval(
			"account.balance >= transaction.withdrawal || (account.overdraftProtection && account.overdraftLimit >= transaction.withdrawal - account.balance)",
			{
				account: {
					balance,
					overdraftProtection,
					overdraftLimit,
				},
				transaction: {
					withdrawal,
				},
			},
		);
		return c.json(data);
	} catch (error) {
		return c.json({ error: `Overdraft test failed: ${error}` }, 500);
	}
});

app.get("/cel/type-analysis", async (c) => {
	try {
		const celery = await getCelery();
		const result = celery.analyzeType("account.balance + transaction.withdrawal", {
			account: {
				balance: 500,
				overdraftProtection: true,
				overdraftLimit: 1000,
			},
			transaction: {
				withdrawal: 700,
			},
		});
		return c.json(result);
	} catch (error) {
		return c.json({ error: `Type analysis failed: ${error}` }, 500);
	}
});

app.get("/cel/type-analysis-unknown", async (c) => {
	try {
		const celery = await getCelery();
		const result = celery.analyzeTypeUnknown(
			"account.balance >= transaction.withdrawal || (account.overdraftProtection && account.overdraftLimit >= transaction.withdrawal - account.balance)",
			["account", "transaction"], // just variable names, no type info
		);
		return c.json(result);
	} catch (error) {
		return c.json({ error: `Unknown type analysis failed: ${error}` }, 500);
	}
});

export default app;
