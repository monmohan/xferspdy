package xfer;

import java.math.BigInteger;
import java.security.MessageDigest;
import java.security.NoSuchAlgorithmException;

/**
 * Implement Adler32 with rolling checksum
 */
public class Hashes {
    public long lv1 = 1, lv2 = 0;
    public static int PRIME = 65521;
    public int st = 0, len = 0, baseIdx = 0;
    public byte[] bytes = null;


    public long adler32(byte[] bytes, int st, int len) {
        this.st = st;
        baseIdx = st;
        this.len = len;
        this.bytes = bytes;
        for (int i = st; i < st + len; i++) {
            lv1 = (lv1 + bytes[i]) % PRIME;
            lv2 = (lv2 + lv1) % PRIME;
        }

        return (lv2 * 65536 + lv1) & 0xffffffffL;

    }

    public long checksum(byte[] bytes) {
        return adler32(bytes, 0, bytes.length);

    }

    public long roll(byte add) {

        lv1 = (lv1 - bytes[st] + add) % PRIME;
        if (lv1 < 0) {
            lv1 = lv1 + PRIME;
        }
        lv2 = (lv2 - (len * bytes[st]) + lv1 - 1) % PRIME;
        if (lv2 < 0) {
            lv2 = lv2 + PRIME;
        }
        bytes[st] = add;
        st++;
        if (st == (baseIdx + len)) {
            st = baseIdx;
        }
        return (lv2 * 65536 + lv1) & 0xffffffffL;

    }

    public long sha() {
        MessageDigest md = null;

        try {
            md = MessageDigest.getInstance("SHA-1");
            md.update(bytes, st, len);
            if (st > baseIdx) {
                md.update(bytes, baseIdx, (st - baseIdx) + 1);
            }
            byte[] dg = md.digest();
            BigInteger bg = new BigInteger(dg);
            md.reset();
            return bg.longValue();
        } catch (NoSuchAlgorithmException e) {
            throw new RuntimeException("no sha");
        }

    }

}
