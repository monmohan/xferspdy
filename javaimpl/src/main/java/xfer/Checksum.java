package xfer;

import java.util.*;

/**
 * Created by singhmo on 4/19/2015.
 */
public class Checksum {
    int blockSz;
    public List<Block> blocks = new ArrayList<Block>();
    public Map<Long, List<Block>> fingerprint = new HashMap<Long, List<Block>>();
    DataSource dataSource = null;

    public Checksum(DataSource dataSource, int blockSz) {
        try {
            this.dataSource = dataSource;
            this.blockSz=blockSz;
        } catch (Exception e) {
            e.printStackTrace();
        }

    }




    public static class Block {
        public Block(int blockSz, long checksum, long SHA, int index) {
            this.blockSz = blockSz;
            this.checksum = checksum;
            this.SHA = SHA;
            this.index = index;
        }

        int blockSz;
        long checksum;
        long SHA;
        int index;

        @Override
        public boolean equals(Object o) {
            if (this == o) return true;
            if (o == null || getClass() != o.getClass()) return false;

            Block block = (Block) o;

            if (SHA != block.SHA) return false;
            if (blockSz != block.blockSz) return false;
            if (checksum != block.checksum) return false;
            if (index != block.index) return false;

            return true;
        }

        @Override
        public int hashCode() {
            int result = blockSz;
            result = 31 * result + (int) (checksum ^ (checksum >>> 32));
            result = 31 * result + (int) (SHA ^ (SHA >>> 32));
            result = 31 * result + index;
            return result;
        }

        @Override
        public String toString() {
            return "Block{" +
                    "blockSz=" + blockSz +
                    ", checksum=" + checksum +
                    ", SHA=" + SHA +
                    ", index=" + index +
                    '}';
        }
    }
}
