package xfer;

import java.io.*;

/**
 * Created by singhmo on 4/23/2015.
 */
public interface DataSource {
    InputStream getStream() throws Exception;



        public static class FileDataSource implements DataSource {
            File source;

            public FileDataSource(File source) {
                this.source = source;
            }

            @Override
            public InputStream getStream() throws Exception{
                return new BufferedInputStream(new FileInputStream(source));
            }
        }

        public static class MemDataSource implements DataSource {
            byte[] source;
            int offset=-1;
            int length=-1;

            public MemDataSource(byte[] source) {
                this.source = source;
                offset=0;
                length=source.length;
            }

            public MemDataSource(byte[] source, int offset,int length ) {
                this.source = source;
                this.length = length;
                this.offset = offset;
            }

            @Override
            public InputStream getStream() throws Exception{
                return new ByteArrayInputStream(source,offset,length);
            }
        }

}
