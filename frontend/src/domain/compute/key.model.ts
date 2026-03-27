// cf. https://zenn.dev/farstep/articles/typescript-branded-types
// branded type によって、別IDの混同などを防止する。 
export type RawPublicKey = string & { readonly __brand: 'RawPublicKey' };
export type SignedCertificate = string & { readonly __brand: 'SignedCertificate' };

export const createRawPublicKey = (key: string): RawPublicKey => {
  if (!key.trim().startsWith('ssh-')) {
    console.warn('Invalid SSH Public Key format detected');
  }
  return key as RawPublicKey;
};

export const createSignedCertificate = (cert: string): SignedCertificate => 
  cert as SignedCertificate;