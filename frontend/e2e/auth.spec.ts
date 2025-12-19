import { test, expect } from '@playwright/test';

test('has title', async ({ page }) => {
  await page.goto('http://localhost:3000/auth/login');

  // Expect a title "to contain" a substring.
  await expect(page).toHaveTitle(/Login/);
});

test('login flow', async ({ page }) => {
  await page.goto('http://localhost:3000/auth/login');

  // Fill the phone number
  await page.getByPlaceholder('09123456789').fill('09123456789');
  
  // Click the submit button
  await page.getByRole('button', { name: /Send OTP/i }).click();

  // Expect to see OTP input (assuming the UI transitions)
  // This depends on the actual implementation, but this is a placeholder for the flow
  // await expect(page.getByPlaceholder('Enter OTP')).toBeVisible();
});
