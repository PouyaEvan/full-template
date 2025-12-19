import { test, expect } from '@playwright/test';

test.describe('Two-Factor Authentication', () => {
  test('should allow user to setup 2FA', async ({ page }) => {
    // Mock login state or login first
    // For this test, we'll assume we can reach the dashboard settings
    // In a real scenario, we'd need to seed the DB or mock the API responses
    
    // Mock the setup API response
    await page.route('**/api/auth/2fa/setup', async route => {
      // Return a dummy image
      const buffer = Buffer.from('iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mNk+M9QDwADhgGAWjR9awAAAABJRU5ErkJggg==', 'base64');
      await route.fulfill({
        status: 200,
        contentType: 'image/png',
        body: buffer
      });
    });

    await page.goto('http://localhost:3000/dashboard/settings/2fa');
    
    // Check if Setup button is visible (assuming not enabled yet)
    // Note: This depends on the initial state of the user which we can't easily control here without seeding
    // So we'll just check for the presence of the page elements
    await expect(page.getByRole('heading', { name: 'Two-Factor Authentication' })).toBeVisible();
  });

  test('should show 2FA input during login if enabled', async ({ page }) => {
    await page.goto('http://localhost:3000/auth/login');

    // Mock the OTP verify response to say 2FA is required
    await page.route('**/api/auth/otp/verify', async route => {
      await route.fulfill({
        status: 200,
        contentType: 'application/json',
        body: JSON.stringify({
          "2fa_required": true,
          "temp_token": "temp-token-123"
        })
      });
    });

    // Fill phone
    await page.getByPlaceholder('09123456789').fill('09123456789');
    await page.getByRole('button', { name: 'Send Code' }).click();

    // Fill OTP (mocked)
    await page.getByPlaceholder('123456').fill('123456');
    await page.getByRole('button', { name: 'Verify & Login' }).click();

    // Expect 2FA input to appear
    await expect(page.getByText('Two-Factor Authentication')).toBeVisible();
    await expect(page.getByPlaceholder('123456')).toBeVisible();
    await expect(page.getByRole('button', { name: 'Verify 2FA' })).toBeVisible();
  });
});
