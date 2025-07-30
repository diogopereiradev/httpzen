import { faker } from '@faker-js/faker';

export const mockRequest = () => {
  return {
    httpVersion: '1.1',
    statusCode: 200,
    responseTime: faker.number.int({ min: 100, max: 500 }),
    method: faker.helpers.arrayElement(['GET', 'POST', 'PUT', 'DELETE']),
    url: faker.internet.url(),
    headers: {
      'Content-Type': 'application/json',
      'User-Agent': faker.internet.userAgent(),
      'Accept-Ch': 'Sec-CH-Prefers-Color-Scheme',
      'Content-Security-Policy': "default-src 'self'; script-src 'none';",
      'X-Xss-Protection': '1; mode=block',
    } as Record<string, string>,
    response: {
      success: true,
      data: {
        message: 'Welcome to Httpzen, your tool to test APIs',
      }
    },
    network: {
      protocol: 'IPv4',
      ip: faker.internet.ipv4(),
      country: faker.location.country(),
      hostname: faker.internet.domainName(),
      city: faker.location.city(),
      decimal: faker.number.int({ min: 1, max: 255 }),
      asn: faker.number.int({ min: 1000, max: 9999 }),
      isp: faker.company.name(),
      coordinates: `${faker.location.latitude()}, ${faker.location.longitude()}`,
    }
  }
}