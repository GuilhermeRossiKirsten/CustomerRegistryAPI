/**
 * Load Test Script for Customer Registry API
 *
 * Testa os endpoints da API com múltiplas operações simultâneas
 * Sem dependências externas - usa apenas built-ins do Node.js
 *
 * Uso: node load-test.js [OPTIONS]
 * Opções:
 *   --users N        Número de usuários simultâneos (default: 10)
 *   --duration N     Duração do teste em segundos (default: 30)
 *   --host URL       URL da API (default: http://localhost:8080)
 */

const { randomUUID } = require("crypto");

// Parse command line arguments
const args = process.argv.slice(2);
let config = {
  users: 10,
  duration: 30,
  host: "http://localhost:8080",
};

for (let i = 0; i < args.length; i += 2) {
  if (args[i] === "--users") config.users = parseInt(args[i + 1]);
  if (args[i] === "--duration") config.duration = parseInt(args[i + 1]);
  if (args[i] === "--host") config.host = args[i + 1];
}

// Função utilitária para fazer requisições HTTP com fetch
async function httpRequest(method, path, data = null) {
  const url = new URL(path, config.host).toString();
  const options = {
    method,
    headers: {
      "Content-Type": "application/json",
    },
    timeout: 5000,
  };

  if (data) {
    options.body = JSON.stringify(data);
  }

  try {
    const response = await fetch(url, options);
    const contentType = response.headers.get("content-type");
    let body = null;

    if (contentType && contentType.includes("application/json")) {
      try {
        body = await response.json();
      } catch (e) {
        body = null;
      }
    }

    return {
      status: response.status,
      data: body,
    };
  } catch (error) {
    throw error;
  }
}

// Métricas
const metrics = {
  requests: 0,
  success: 0,
  errors: 0,
  totalTime: 0,
  minTime: Infinity,
  maxTime: 0,
  statusCodes: {},
  errors: [],
};

// Dados de teste
const testData = {
  customers: [],
  documents: [],
};

// Função para gerar dados de cliente
function generateCustomer() {
  const id = randomUUID();
  const document = `DOC-${Math.random()
    .toString(36)
    .substr(2, 9)
    .toUpperCase()}`;
  const risks = ["LOW", "MEDIUM", "HIGH"];
  const statuses = ["ACTIVE", "INACTIVE", "UNDER_REVIEW"];

  return {
    document,
    name: `Customer ${document}`,
    score: Math.floor(Math.random() * 1000),
    risk_level: risks[Math.floor(Math.random() * risks.length)],
    income_range: `${1000 + Math.floor(Math.random() * 9000)}-${
      10000 + Math.floor(Math.random() * 90000)
    }`,
    status: statuses[Math.floor(Math.random() * statuses.length)],
  };
}

// Função para registrar métrica
function recordMetric(statusCode, time) {
  metrics.requests++;
  metrics.totalTime += time;
  metrics.minTime = Math.min(metrics.minTime, time);
  metrics.maxTime = Math.max(metrics.maxTime, time);

  if (statusCode < 400) {
    metrics.success++;
  } else {
    metrics.errors++;
  }

  metrics.statusCodes[statusCode] = (metrics.statusCodes[statusCode] || 0) + 1;
}

// Operação: Criar cliente
async function createCustomer() {
  const customer = generateCustomer();
  const startTime = Date.now();

  try {
    const response = await httpRequest("POST", "/customers", customer);
    const duration = Date.now() - startTime;
    recordMetric(response.status, duration);

    if (response.status === 201 && response.data) {
      testData.customers.push(response.data);
      testData.documents.push(customer.document);
    }

    return response.data;
  } catch (error) {
    const duration = Date.now() - startTime;
    recordMetric(0, duration);
  }
}

// Operação: Listar clientes
async function listCustomers() {
  const startTime = Date.now();

  try {
    const limit = 20;
    const offset = Math.floor(Math.random() * 100);
    const response = await httpRequest(
      "GET",
      `/customers?limit=${limit}&offset=${offset}`
    );
    const duration = Date.now() - startTime;
    recordMetric(response.status, duration);
  } catch (error) {
    const duration = Date.now() - startTime;
    recordMetric(0, duration);
  }
}

// Operação: Buscar cliente por ID
async function getCustomerById() {
  if (testData.customers.length === 0) return;

  const customer =
    testData.customers[Math.floor(Math.random() * testData.customers.length)];
  const startTime = Date.now();

  try {
    const response = await httpRequest("GET", `/customers/${customer.id}`);
    const duration = Date.now() - startTime;
    recordMetric(response.status, duration);
  } catch (error) {
    const duration = Date.now() - startTime;
    recordMetric(0, duration);
  }
}

// Operação: Buscar cliente por documento
async function getCustomerByDocument() {
  if (testData.documents.length === 0) return;

  const document =
    testData.documents[Math.floor(Math.random() * testData.documents.length)];
  const startTime = Date.now();

  try {
    const response = await httpRequest(
      "GET",
      `/customers/document/${document}`
    );
    const duration = Date.now() - startTime;
    recordMetric(response.status, duration);
  } catch (error) {
    const duration = Date.now() - startTime;
    recordMetric(0, duration);
  }
}

// Operação: Atualizar status
async function updateStatus() {
  if (testData.customers.length === 0) return;

  const customer =
    testData.customers[Math.floor(Math.random() * testData.customers.length)];
  const statuses = ["ACTIVE", "INACTIVE", "UNDER_REVIEW"];
  const newStatus = statuses[Math.floor(Math.random() * statuses.length)];
  const startTime = Date.now();

  try {
    const response = await httpRequest(
      "PATCH",
      `/customers/${customer.id}/status`,
      {
        status: newStatus,
      }
    );
    const duration = Date.now() - startTime;
    recordMetric(response.status, duration);
  } catch (error) {
    const duration = Date.now() - startTime;
    recordMetric(0, duration);
  }
}

// Worker que executa operações aleatórias
async function worker() {
  const operations = [
    createCustomer,
    listCustomers,
    getCustomerById,
    getCustomerByDocument,
    updateStatus,
  ];

  while (!shouldStop) {
    const operation = operations[Math.floor(Math.random() * operations.length)];
    try {
      await operation();
    } catch (error) {
      console.error(`Worker erro: ${error.message}`);
    }

    // Pequeno delay entre requisições
    await new Promise((resolve) => setTimeout(resolve, 100));
  }
}

// Flag para parar o teste
let shouldStop = false;

// Função principal
async function runLoadTest() {
  console.log("\n╔════════════════════════════════════════╗");
  console.log("║  Customer Registry API - Load Test     ║");
  console.log("╚════════════════════════════════════════╝\n");

  console.log(`📊 Configuração:`);
  console.log(`   • Usuários simultâneos: ${config.users}`);
  console.log(`   • Duração: ${config.duration}s`);
  console.log(`   • Host: ${config.host}\n`);

  console.log("⏱️  Iniciando teste...\n");

  // Criar workers
  const workers = [];
  for (let i = 0; i < config.users; i++) {
    workers.push(worker());
  }

  // Mostrar progresso
  const startTime = Date.now();
  const progressInterval = setInterval(() => {
    const elapsed = Math.floor((Date.now() - startTime) / 1000);
    const avgTime =
      metrics.requests > 0
        ? (metrics.totalTime / metrics.requests).toFixed(2)
        : 0;
    const rps = (metrics.requests / elapsed).toFixed(2);

    console.log(
      `[${elapsed}s/${config.duration}s] Requisições: ${metrics.requests} | ✓ ${metrics.success} | ✗ ${metrics.errors} | RPS: ${rps} | Média: ${avgTime}ms`
    );
  }, 5000);

  // Aguardar duração do teste
  await new Promise((resolve) => setTimeout(resolve, config.duration * 1000));

  // Parar teste
  shouldStop = true;
  await Promise.all(workers);
  clearInterval(progressInterval);

  // Exibir resultados
  const totalTime = Date.now() - startTime;
  const avgTime = metrics.totalTime / metrics.requests;
  const rps = (metrics.requests / (totalTime / 1000)).toFixed(2);
  const successRate = ((metrics.success / metrics.requests) * 100).toFixed(2);

  console.log("\n╔════════════════════════════════════════╗");
  console.log("║         RESULTADOS DO TESTE            ║");
  console.log("╚════════════════════════════════════════╝\n");

  console.log(`⏱️  Tempo total: ${(totalTime / 1000).toFixed(2)}s`);
  console.log(`📊 Total de requisições: ${metrics.requests}`);
  console.log(`✓ Sucesso: ${metrics.success} (${successRate}%)`);
  console.log(`✗ Erros: ${metrics.errors}`);
  console.log(`\n⚡ Performance:`);
  console.log(`   • RPS (Requisições/s): ${rps}`);
  console.log(`   • Tempo médio: ${avgTime.toFixed(2)}ms`);
  console.log(`   • Tempo mínimo: ${metrics.minTime}ms`);
  console.log(`   • Tempo máximo: ${metrics.maxTime}ms`);

  console.log(`\n📈 Distribuição de Status HTTP:`);
  Object.entries(metrics.statusCodes)
    .sort((a, b) => b[1] - a[1])
    .forEach(([status, count]) => {
      const percentage = ((count / metrics.requests) * 100).toFixed(2);
      console.log(`   • ${status}: ${count} (${percentage}%)`);
    });

  console.log(`\n✅ Teste finalizado!\n`);
}

// Tratamento de erros
process.on("unhandledRejection", (err) => {
  console.error("❌ Erro não tratado:", err);
});

// Executar teste
runLoadTest().catch((err) => {
  console.error("❌ Erro ao executar teste:", err);
  process.exit(1);
});
