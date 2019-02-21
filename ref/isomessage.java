

package iso;

import util.Tool;

import java.io.*;
import java.net.Socket;
import java.util.HashMap;

/**
 * @author iyus
 * version 2.4 - Add setNumber
 *
 */
public class IsoMessage {
    public static final String DLM = "|";

    protected String mti;
    protected HashMap<Integer, FieldDef> fields;
    protected String[] bitmapValues = new String[128];
    protected int[] bitmaps = new int[16];
    private int bits[] = {0, 1, 128, 64, 32, 16, 8, 4, 2};
    protected int etx = -1;
    protected boolean debug = false;

    @Deprecated
    public IsoMessage(HashMap<Integer, FieldDef> fields) {
        for (int i = 0; i < 16; i++)
            bitmaps[i] = 0;
        for (int i = 0; i < 128; i++)
            bitmapValues[i] = "";

        this.fields = fields;
    }

    public IsoMessage(FieldDef fieldDefs[]) {
        for (int i = 0; i < 16; i++)
            bitmaps[i] = 0;
        for (int i = 0; i < 128; i++)
            bitmapValues[i] = "";

        fields = new HashMap<Integer, FieldDef>();
        for (FieldDef fieldDef : fieldDefs)
            fields.put(fieldDef.bit, fieldDef);
    }

    public IsoMessage() {
        for (int i = 0; i < 16; i++)
            bitmaps[i] = 0;
        for (int i = 0; i < 128; i++)
            bitmapValues[i] = "";
        fields = new HashMap<Integer, FieldDef>();
    }

    public void encode(byte[] buffer, boolean hasLength) throws Exception {
        // tambahan 2016-04-29
        encode(new String(buffer), hasLength);
    }

    public void encode(byte[] buffer) throws Exception {
        encode(new String(buffer));
    }

    public void encode(String source) throws Exception {
        encode(source, false);
    }

    public void encode(String source, boolean hasLength) throws Exception {
        int offset = hasLength ? 4 : 0; // hanya berlaku untuk length 4
        if (source.length() < 20 + offset)
            throw new Exception("Invalid iso string: " + source);

        this.mti = source.substring(offset, offset + 4);
        offset += 4;
        String bitmapHex = this.buildBitmap(source, offset);
        offset += bitmapHex.length();
        buildValues(bitmapHex, source, offset);
    }

    public void execute(String address) throws Exception {
        String[] pair = address.split(":");
        if (pair.length != 2)
            throw new Exception("Invalid address: " + address);

        String host = pair[0];
        int port = Integer.parseInt(pair[1]);

        execute(host, port);
    }

    public void execute(String host, int port) throws Exception {
        //TODO cek lagi apakah pattern nya udah bener ?
        DataInputStream inputStream;
        DataOutputStream outputStream;
        Object serverInLock = new Object();
        Object serverOutLock = new Object();
        //String raw = this.decode(false);
        String raw = this.decode();

        Socket socket = new Socket(host, port);

        synchronized (serverOutLock) {
            // send message
            outputStream = new DataOutputStream(new BufferedOutputStream(socket.getOutputStream(), 2048));
            outputStream.writeBytes(raw);
            outputStream.write(etx);
            outputStream.flush();
        }

        synchronized (serverInLock) {
            // receive response
            inputStream = new DataInputStream(new BufferedInputStream(socket.getInputStream()));

            if (hasEtx()) {
                int cell, k = 0;
                StringBuffer sb = new StringBuffer();
                while ((cell = inputStream.read()) != etx) {
                    sb.append((char) cell);
                    k += 1;
                    if (k > 2048) {
                        System.out.println("Kepanjangan ...");
                        break;
                    }
                }
                raw = sb.toString() ;

                if (debug)
                    System.out.println(">" + raw);
            } else {
                int len = getLen(inputStream);

                byte[] buffer = new byte[len];
                inputStream.readFully(buffer);
                raw = new String(buffer);
            }
        }
        this.encode(raw, false);

        socket.close();
    }

    protected int getLen(DataInputStream inputStream) throws IOException {
        byte[] lenbuf = new byte[4];
        inputStream.readFully(lenbuf);
        return Integer.parseInt(new String(lenbuf));
    }

    public void write(DataOutputStream outputStream) {
        try {
            outputStream.writeBytes(this.decode());
            outputStream.flush();
        } catch (IOException e) {
        }
    }

    protected void buildValues(String bitmapHex, String source, int offset) throws Exception {
        // empty bitmapValues
        for (int i = 0; i < 128; i++)
            bitmapValues[i] = "";

        // flag yang akan diisi
        for (int i = 0; i < 16; i++)
            for (int j = 1; j < 9; j++)
                if ((bitmaps[i] & bits[j]) == bits[j])
                    if (j == 1)
                        bitmapValues[(i + 1) * 8] = "X";
                    else if (j != 2 || i != 0)
                        bitmapValues[i * 8 + j - 1] = "X";

        boolean showbit = false; // for debug only
        if (showbit) {
            for (int bit = 2; bit < 128; bit++)
                if (bitmapValues[bit].equals("X"))
                    System.out.println("Bit " + bit);
            throw new Exception("SENGAJA");
        }

        for (int bit = 2; bit < 128; bit++) {
            Integer key = new Integer(bit);

            if (bitmapValues[bit].equals("X") && fields.containsKey(key)) {
                FieldDef field = fields.get(key);

                if (field.isLL()) {
                    if (debug)
                        System.out.println("Will process Bit LL " + bit);

                    int valsize = Integer.parseInt(source.substring(offset, offset + 2));
                    if (valsize > field.length)
                        throw new Exception("Length exceed defined in bit " + bit + ", size: " + valsize + ", field.length: " + field.length);
                    bitmapValues[bit] = source.substring(offset, offset + 2) +
                            source.substring(offset + 2, offset + 2 + valsize);
                    offset += valsize + 2;
                } else if (field.isLLL()) {
                    if (debug)
                        System.out.println("Will process Bit LLL " + bit);

                    int valsize = Integer.parseInt(source.substring(offset, offset + 3));
                    if (valsize > field.length)
                        throw new Exception("Length exceed defined in bit " + bit);
                    bitmapValues[bit] = source.substring(offset, offset + 3) +
                            source.substring(offset + 3, offset + 3 + valsize);
                    offset += valsize + 3;
                } else {
                    if (debug)
                        System.out.println("Will process Bit " + bit);

                    bitmapValues[bit] = source.substring(offset, offset + field.length);
                    offset += field.length;
                }

                if (debug)
                    System.out.println("    Bit " + bit + " : [" + bitmapValues[bit] + "]");
            }
        }
    }

    public String decode()
    {
        return decode(hasEtx() ? false : true);
    }

    public String decode(boolean hasLength) {
        StringBuilder result = new StringBuilder();
        result.append(this.mti);
        result.append(this.buildBitmap());

        for (int i = 0; i < 128; i++)
            if (bitmapValues[i] != null && bitmapValues[i].length() > 0)
                result.append(bitmapValues[i]);

        String isostr = result.toString();

        if (hasLength) {
            String len = getStringLen(isostr);
            isostr = len + isostr;
        }

        return isostr;
    }

    protected String getStringLen(String isostr) {
        String len = String.format("%4d", isostr.length());
        len = len.replace(' ', '0');
        return len;
    }

    public void assign(IsoMessage source) throws Exception {
        for (int i = 2; i < 128; i++) {
            String val = source.get(i);
            this.set(i, val);
        }
        this.mti = source.mti;
    }

    protected String buildBitmap() {
        String bitmapHex = "";
        for (int c = 0; c < 16; c++) {
            String tm = Integer.toHexString(bitmaps[c]);
            if (tm.length() < 2)
                tm = "0" + tm;
            bitmapHex += tm;

            if ((bitmaps[0] & 128) != 128 && (c == 7)) break;
        }
        //return bitmapHex;
        return bitmapHex.toUpperCase();
    }

    protected String buildBitmap(String source, int offset) {
        int top = offset + 32;
        int mid = offset + 14;
        String bitmapHex = "";
        int k = 0;
        String s = source.substring(offset, offset + 2);
        int first = Integer.parseInt(s, 16);
        boolean noext = (first & 128) != 128;

        while (offset < top) {
            String tmp = source.substring(offset, offset + 2);
            bitmapHex += tmp;
            bitmaps[k] = Integer.parseInt(tmp, 16);

            if (noext && (offset == mid))
                break;

            k++;
            offset += 2;
        }
        //return bitmapHex;
        return bitmapHex.toUpperCase();
    }

    public void set(int index, Object value) throws Exception {
        if (value == null) return;
        if (value.toString().equals("")) return;

        Integer key = new Integer(index);
        if (fields.containsKey(key)) {
            FieldDef def = fields.get(key);
            bitmapValues[index] = def.format(value);

            if (index > 64)
                bitmaps[0] = bitmaps[0] | bits[2];
            int pos = (index % 8) == 0 ? (index / 8) - 1 : index / 8;
            bitmaps[pos] = bitmaps[pos] | bits[(index % 8) + 1];
        }
    }

    public IsoMessage setString(int index, String value) {
        try {
            set(index, value);
        } catch (Exception e) {
            e.printStackTrace();
        }
        return this;
    }

    public IsoMessage setNumber(int index, String value) {
        try {
            value = Tool.fix(value);
            set(index, value);
        } catch (Exception e) {
            e.printStackTrace();
        }
        return this;
    }

    public static void validate(boolean valid, String message) throws Exception {
        if (!valid)
            throw new Exception(message);
    }

    public boolean hasIndex(int index) {
        Integer key = new Integer(index);
        return fields.containsKey(key) && bitmapValues[index].length() > 0;
    }

    public String get(int index) {
        Integer key = new Integer(index);
        if (fields.containsKey(key) && bitmapValues[index].length() > 0) {
            FieldDef f = fields.get(key);
            String s = bitmapValues[index];

            if (f.isLL())
                return s.substring(2);
            else if (f.isLLL())
                return s.substring(3);
            else
                return s;
        } else
            return null;
    }

    public String toString() {
        int bit = 0;
        try {
            StringBuilder result = new StringBuilder();
            result.append("MTI: [" + mti + "]\n");
            for (bit = 2; bit < 128; bit++) {
                String x = get(bit);
                if (x != null)
                    result.append(String.format("%3d: [%s]\n", bit, x));
            }

            return result.toString();
        } catch (Exception e) {
            return String.format("Error at bit %d\nReason:%s\n", bit, e.getMessage());
        }
    }

    public String getMti() {
        return mti;
    }

    public IsoMessage setMti(String mti) {
        this.mti = mti;
        return this;
    }

    public String getString(int bit) {
        String s = get(bit);
        if (s == null)
            return "";
        return s.trim();
    }

    public long optLong(int bit, long def) {
        try {
            return Long.parseLong(get(bit));
        } catch (NumberFormatException e) {
            return def;
        }
    }

    public int optInt(int bit, int def) {
        try {
            return Integer.parseInt(get(bit));
        } catch (NumberFormatException e) {
            return def;
        }
    }

    public String optString(int bit, String def) {
        if (hasIndex(bit))
            return getString(bit);
        else
            return def;
    }

    public void setField(FieldDef fieldDef) {
        fields.put(fieldDef.bit, fieldDef);
    }

    public void set(int bit, Object... args) throws Exception {
        set(bit, combine(args));
    }

    public static String combine(Object... args) {
        StringBuilder sb = new StringBuilder();
        for (int i = 0; i < args.length; i++) {
            String s = args[i] == null ? "" : args[i].toString();
            if (i == args.length - 1)
                sb.append(s);
            else
                sb.append(s + DLM);
        }
        return sb.toString();
    }

    public String[] split(int bit) {
        return getString(bit).split("\\" + DLM);
    }

    public void setEtx(int etx) {
        this.etx = etx;
    }

    public boolean hasEtx() {
        return this.etx > -1;
    }

    public int getExt() {
        return this.etx;
    }

    public void setDebug(boolean value) {
        this.debug = value;
    }

    public String getBillerCode() {
        return getString(100);
    }

    public String getProcessingCode() {
        return getString(3);
    }

    public String getReturnCode() {
        return getString(39);
    }
}

