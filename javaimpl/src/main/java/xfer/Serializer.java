package xfer;

import java.io.*;
import java.util.ArrayList;
import java.util.List;

/**
 * Created by singhmo on 4/21/2015.
 */
public class Serializer {

    public static Checksum toChecksum(DataSource ds, int blockSz) throws Exception {
        Checksum cs = new Checksum(ds, blockSz);
        InputStream bis = ds.getStream();
        byte[] chunk = new byte[blockSz];
        int bytes = bis.read(chunk);
        int index = 0;
        Hashes hashes=null;
        while (bytes != -1) {
            hashes=new Hashes();
            long csum = hashes.adler32(chunk, 0, bytes);
            long sha = hashes.sha();
            Checksum.Block block = new Checksum.Block(bytes, csum, sha, index);
            cs.blocks.add(block);
            List<Checksum.Block> blk = cs.fingerprint.get(csum);
            if (blk == null) {
                blk = new ArrayList<Checksum.Block>(10);
                cs.fingerprint.put(csum, blk);
            }
            blk.add(block);
            index++;
            bytes = bis.read(chunk);
        }
        return cs;
    }

    public static void fromChecksum(Checksum cs, FileOutputStream fos) throws Exception {
        InputStream fin = cs.dataSource.getStream();
        byte[] buf = new byte[cs.blockSz];
        for (Checksum.Block block : cs.blocks) {
            int bytes = fin.read(buf);
            fos.write(buf, 0, bytes);
        }
        fin.close();
        fos.close();
    }

    public static Diff diff(Checksum source, DataSource target) throws Exception {
        Diff diff = new Diff();

        InputStream bis = target.getStream();
        byte[] chunk = new byte[source.blockSz], single = new byte[1];
        int rangeStart = 0, rangeEnd = 0, bytes = -1;
        Hashes cs = null;
        Checksum.Block blk = null;
        boolean roll = false;
        while (bis.available() != 0) {
            long csum = -1;
            if (roll) {
                bytes = bis.read(single);
                csum = cs.roll(single[0]);
            } else {
                cs = new Hashes();
                bytes = bis.read(chunk);
                csum = cs.adler32(chunk, 0, bytes);
            }

            if ((blk = blockMatched(source, cs, csum)) != null) {
                if (rangeStart != rangeEnd) {
                    Diff.ByteRange range = new Diff.ByteRange(rangeStart, rangeEnd);
                    diff.diffData.add(range);
                }
                diff.diffData.add(blk);
                rangeStart += bytes;
                rangeEnd = rangeStart;
                roll = false;
            } else {
                rangeEnd++;
                roll = true;
            }
        }
        return diff;

    }

    private static Checksum.Block blockMatched(Checksum source, Hashes cs,long csum) {
        List<Checksum.Block> existing = source.fingerprint.get(csum);
        Checksum.Block matched = null;
        if (existing != null) {
            long sha = cs.sha();
            for (
                Checksum.Block block : existing) {
                if (sha == block.SHA) {
                    matched = block;
                    break;
                }
            }
            if (matched == null) {
                System.out.println("32 bit checksum matched but SHA hash " +
                        "match was not found");
            }
        }
        return matched;
    }


}
